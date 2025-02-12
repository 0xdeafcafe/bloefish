package cher

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/0xdeafcafe/bloefish/libraries/stacktrace"
	"github.com/pkg/errors"
)

// errors that are expected to be common across services
const (
	BadRequest        = "bad_request"
	Unauthorized      = "unauthorized"
	AccessDenied      = "access_denied"
	NotFound          = "not_found"
	RouteNotFound     = "route_not_found"
	MethodNotAllowed  = "method_not_allowed"
	Unknown           = "unknown"
	NoLongerSupported = "no_longer_supported"
	TooManyRequests   = "too_many_requests"
	ContextCanceled   = "context_canceled"
	EOF               = "eof"
	UnexpectedEOF     = "unexpected_eof"
	RequestTimeout    = "request_timeout"
	ThirdPartyTimeout = "third_party_timeout"

	CoercionError = "unable_to_coerce_error"
)

// E implements the official CHER structure
type E struct {
	Code    string             `json:"code"`
	Meta    M                  `json:"meta,omitempty"`
	stack   []stacktrace.Frame `json:"-"`
	Reasons []E                `json:"reasons,omitempty"`
}

// New returns a new E structure with code, meta, and optional reasons.
func New(code string, meta M, reasons ...E) E {
	// For errors with only one reason, we can merge stack traces so we can actually see where the error originally came from
	stack := stacktrace.GetCallerFrames(2)
	if len(reasons) == 1 {
		rootStack := reasons[0].stack
		stack = stacktrace.MergeStacks(rootStack, stack)
	}
	return E{
		Code:    code,
		Meta:    meta,
		stack:   stack,
		Reasons: reasons,
	}
}

// Errorf returns a new E structure, with a message formatted by fmt.
func Errorf(code string, meta M, format string, args ...any) E {
	meta["message"] = fmt.Sprintf(format, args...)

	return E{
		Code: code,
		Meta: meta,
	}
}

// StatusCode returns the HTTP Status Code associated with the
// current error code.
// Defaults to 400 Bad Request because if something's explicitly
// handled with Cher, it is considered "by design" and not
// worthy of a 500, which will alert.
func (e E) StatusCode() int {
	switch e.Code {
	case Unauthorized:
		return http.StatusUnauthorized

	case AccessDenied:
		return http.StatusForbidden

	case NotFound, RouteNotFound:
		return http.StatusNotFound

	case MethodNotAllowed:
		return http.StatusMethodNotAllowed

	case NoLongerSupported:
		return http.StatusGone

	case TooManyRequests:
		return http.StatusTooManyRequests

	case Unknown, CoercionError, RequestTimeout:
		return http.StatusInternalServerError
	}

	return http.StatusBadRequest
}

// Error implements the error interface.
func (e E) Error() string {
	return e.Code
}

// Serialize returns a json representation of the CHER structure
func (e E) Serialize() string {
	output, err := json.Marshal(e)
	if err != nil {
		return ""
	}

	return string(output)
}

func (e E) GetStack() []stacktrace.Frame {
	return e.stack
}

// M it an alias type for map[string]any
type M map[string]any

// Coerce attempts to coerce a CHER out of any object.
// - `E` types are just returned as-is
// - strings are taken as the Code for an E object
// - bytes are unmarshaled from JSON to an E object
// - types implementing the `error` interface to an E object with the error as a reason
func Coerce(val any) E {
	switch val := val.(type) {
	case E:
		return val

	case string:
		return E{Code: val}

	case []byte:
		var cerr E

		err := json.Unmarshal(val, &cerr)
		if err != nil {
			return E{
				Code: CoercionError,
				Meta: M{
					"message": err.Error(),
				},
			}
		}

		return cerr

	case error:
		val = errors.Cause(val)

		return E{
			Code: Unknown,
			Meta: M{
				"message": val.Error(),
			},
		}
	}

	return E{Code: CoercionError}
}

func (e E) Value() (driver.Value, error) {
	return json.Marshal(e)
}

// WrapIfNotCher will not wrap the error if it is any cher except unknown.
//
// If you know the codes you could get back, you should use WrapIfNotCherCodes.
func WrapIfNotCher(err error, msg string) error {
	if err == nil {
		return nil
	}

	var cErr E
	if errors.As(err, &cErr) {
		if cErr.Code == Unknown {
			return errors.Wrap(err, msg)
		}

		return cErr
	}

	return errors.Wrap(err, msg)
}

// WrapIfNotCherCodes will wrap an error unless it is a cher with specific codes.
func WrapIfNotCherCodes(err error, msg string, codes []string) error {
	return WrapIfNotCherCode(err, msg, codes...)
}

func WrapIfNotCherCode(err error, msg string, codes ...string) error {
	var cErr E
	if errors.As(err, &cErr) && slices.Contains(codes, cErr.Code) {
		return cErr
	}

	return errors.Wrap(err, msg)
}

func AsCherWithCode(err error, codes ...string) (cErr E, ok bool) {
	return cErr, errors.As(err, &cErr) && slices.Contains(codes, cErr.Code)
}
