package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/libraries/contexts"
	"github.com/0xdeafcafe/bloefish/libraries/errfuncs"
	"github.com/0xdeafcafe/bloefish/libraries/merr"
	"github.com/0xdeafcafe/bloefish/libraries/mlog"
	"github.com/0xdeafcafe/bloefish/libraries/slicefuncs"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/semconv/v1.13.0/httpconv"
	"go.opentelemetry.io/otel/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type responseWriter struct {
	http.ResponseWriter

	Status int
	Bytes  int64
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.Status == 0 {
		rw.Status = code
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(bytes []byte) (int, error) {
	if rw.Status == 0 {
		rw.Status = http.StatusOK
		rw.WriteHeader(http.StatusOK)
	}

	rw.Bytes += int64(len(bytes))

	return rw.ResponseWriter.Write(bytes)
}

// Telemetry returns a middleware handler that wraps subsequent middleware/handlers and logs
// request information AFTER the request has completed. It also injects a request-scoped
// logger on the context which can be set, read and updated using clog lib
//
// Included fields:
//   - Request ID                (request_id)
//   - HTTP Method               (http_method)
//   - HTTP Path                 (http_path)
//   - HTTP Protocol Version     (http_proto)
//   - Remote Address            (http_remote_addr)
//   - User Agent Header         (http_user_agent)
//   - Referer Header            (http_referer)
//   - Duration with unit        (http_duration)
//   - Duration in microseconds  (http_duration_us)
//   - HTTP Status Code          (http_status)
//   - Response in bytes         (http_response_bytes)
//   - Client Version header     (http_client_version)
//   - User Agent header         (http_user_agent)
func Telemetry(log *logrus.Entry) func(http.Handler) http.Handler {
	propagator := otel.GetTextMapPropagator()
	tracer := otel.Tracer(
		"github.com/0xdeafcafe/bloefish/libraries/crpc",
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := tracer.Start(
				propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header)),
				"crpc", // Placeholder, this will be updates below
				trace.WithSpanKind(trace.SpanKindServer),
			)
			defer span.End()

			// inject the span context into the response headers
			propagator.Inject(ctx, propagation.HeaderCarrier(w.Header()))

			// update the span name with the service name and endpoint
			svcCtx := contexts.MustGetServiceInfo(ctx)
			span.SetName(fmt.Sprintf("crpc (%s%s)", svcCtx.ServiceHTTPName, r.URL.Path))

			// create a mutable logger instance which will persist for the request
			// inject pointer to the logger into the request context
			ctx = clog.Set(ctx, log)
			r = r.WithContext(ctx)

			// panics inside handlers will be logged to standard before propagation
			defer clog.HandlePanic(ctx, true)

			// set useful log fields for the request
			clog.SetFields(ctx, clog.Fields{
				"http_remote_addr":    r.RemoteAddr,
				"http_user_agent":     r.UserAgent(),
				"http_client_version": r.Header.Get("Infra-Client-Version"),
				"http_path":           r.URL.Path,
				"http_method":         r.Method,
				"http_proto":          r.Proto,
				"http_referer":        r.Referer(),
			})

			// wrap given response writer with one that tracks status code/bytes written
			res := &responseWriter{ResponseWriter: w}

			tStart := time.Now()
			next.ServeHTTP(res, r)
			tEnd := time.Now()

			// set useful log fields from the response
			clog.SetFields(ctx, clog.Fields{
				"http_duration":       tEnd.Sub(tStart).String(),
				"http_duration_us":    int64(tEnd.Sub(tStart) / time.Microsecond),
				"http_status":         res.Status,
				"http_response_bytes": res.Bytes,
			})

			// set the span attributes and status
			span.SetAttributes(semconv.HTTPResponseStatusCode(res.Status))
			span.SetStatus(httpconv.ServerStatus(res.Status))

			// get the error if one is set on the log entry
			err := getError(clog.Get(ctx))
			if err == nil {
				mlog.Info(ctx, merr.New(ctx, "request_completed", nil))
				return
			}

			// if appropriate, log the error at the appropriate level
			var fn func(context.Context, merr.Merrer)
			switch clog.DetermineLevel(err, clog.TimeoutsAsErrors(ctx)) {
			case
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel:
				fn = mlog.Error
			case logrus.WarnLevel:
				fn = mlog.Warn
			case
				logrus.InfoLevel,
				logrus.DebugLevel,
				logrus.TraceLevel:
				fn = mlog.Info
			}

			if mErr, ok := errfuncs.As[merr.E](err); ok {
				fn(ctx, mErr)
			} else if cErr, ok := errfuncs.As[cher.E](err); ok {
				reasons := slicefuncs.Map(cErr.Reasons, func(r cher.E) error { return r })

				// If the cher error has no reasons, add the cher error itself
				if len(reasons) == 0 {
					reasons = append(reasons, cErr)
				}

				fn(ctx, merr.New(ctx, merr.Code(cErr.Code), merr.M(cErr.Meta), reasons...))
			} else {
				fn(ctx, merr.New(ctx, "unexpected_request_failure", nil, err))
			}
		})
	}
}

// getError returns the error if one is set on the log entry
func getError(l *logrus.Entry) error {
	if erri, ok := l.Data[logrus.ErrorKey]; ok {
		if err, ok := erri.(error); ok {
			return err
		}
	}

	return nil
}
