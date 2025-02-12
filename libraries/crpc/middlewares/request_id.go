package middlewares

import (
	"context"
	"net/http"

	"github.com/0xdeafcafe/bloefish/libraries/contexts"
	"github.com/0xdeafcafe/bloefish/libraries/ksuid"
)

// GetRequestID returns the request id embedded within a HTTP request's context,
// or an empty string if no request id has been established.
func GetRequestID(r *http.Request) string {
	return contexts.GetRequestID(r.Context())
}

// SetRequestID sets the request id within a HTTP requests's context,
// it will overwrite any existing request id.
func SetRequestID(r *http.Request, requestID string) *http.Request {
	return r.WithContext(contexts.SetRequestID(r.Context(), requestID))
}

// GetOrSetRequestID will return any request ID found in the context, and
// if one does not exist, set and return.
func GetOrSetRequestID(ctx context.Context) (context.Context, string) {
	requestID := contexts.GetRequestID(ctx)
	if requestID != "" {
		return ctx, requestID
	}

	requestID = ksuid.Generate(ctx, "req").String()

	ctx = contexts.SetRequestID(ctx, requestID)
	return ctx, requestID
}

// RequestID either generates a new request id or acquires an existing request id
// from the http requests headers, and embeds into the requests context and
// response headers.
func RequestID(next http.Handler) http.Handler {
	const requestIDHeader = "Request-Id"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(requestIDHeader)
		if requestID == "" {
			// if no Request-ID is passed, generate/originate a new one
			requestID = ksuid.Generate(r.Context(), "req").String()
		}

		w.Header().Set(requestIDHeader, requestID)

		r = SetRequestID(r, requestID)

		next.ServeHTTP(w, r)
	})
}
