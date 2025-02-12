package crpc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/libraries/errfuncs"
	"github.com/0xdeafcafe/bloefish/libraries/merr"
	"github.com/0xdeafcafe/bloefish/libraries/mlog"
	"github.com/0xdeafcafe/bloefish/libraries/slicefuncs"

	"github.com/xeipuuv/gojsonschema"
)

// Logger inherits the context logger and reports RPC request success/failure.
func Logger() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(res http.ResponseWriter, req *Request) error {
			ctx := req.Context()

			// Add fields to the request context scoped logger, if one exists
			clog.SetFields(ctx, clog.Fields{
				"rpc_version": req.Version,
				"rpc_method":  req.Method,
			})

			tStart := time.Now()
			err := next(res, req)
			tEnd := time.Now()

			clog.SetFields(ctx, clog.Fields{
				"rpc_duration":    tEnd.Sub(tStart).String(),
				"rpc_duration_us": int64(tEnd.Sub(tStart) / time.Microsecond),
			})

			if err == nil {
				return nil
			}

			// rewrite common errors to internal error standard
			switch {
			case errors.Is(err, io.EOF):
				err = cher.New(cher.EOF, nil)
			case errors.Is(err, io.ErrUnexpectedEOF):
				err = cher.New(cher.UnexpectedEOF, nil)
			case errors.Is(err, context.Canceled):
				err = cher.New(cher.ContextCanceled, nil)
			}

			clog.SetError(ctx, err)

			return err
		}
	}
}

// Validate buffers the JSON body and applies a JSON Schema validation.
func Validate(schema *gojsonschema.Schema) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(res http.ResponseWriter, req *Request) error {
			ctx := req.Context()

			body, err := io.ReadAll(req.Body)
			if err != nil {
				if netErr, ok := errfuncs.As[net.Error](err); ok {
					mlog.Warn(ctx, merr.New(ctx, "request_body_read_network_error", nil, netErr))
					return io.ErrUnexpectedEOF
				}

				return merr.New(ctx, "request_body_read_failed", nil, err)
			}

			result, err := schema.Validate(gojsonschema.NewBytesLoader(body))
			if err != nil {
				if errors.Is(err, io.EOF) {
					return cher.New("invalid_json", nil)
				}

				return merr.New(ctx, "request_body_validation_failed", nil, err)
			}

			err = CoerceJSONSchemaError(result)
			if err != nil {
				return err
			}

			req.Body = io.NopCloser(bytes.NewReader(body))
			return next(res, req)
		}
	}
}

func CoerceJSONSchemaError(result *gojsonschema.Result) error {
	if result.Valid() {
		return nil
	}

	return cher.New(cher.BadRequest, nil, slicefuncs.Map(result.Errors(), func(err gojsonschema.ResultError) cher.E {
		return cher.E{
			Code: "schema_failure",
			Meta: cher.M{
				"field":   err.Field(),
				"type":    err.Type(),
				"message": err.Description(),
			},
		}
	})...)
}
