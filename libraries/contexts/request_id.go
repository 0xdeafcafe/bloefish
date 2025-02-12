package contexts

import "context"

type requestIDContextKey string

const RequestIDKey requestIDContextKey = "Request-ID"

// GetRequestID returns the request id embedded within a context, or an empty string if no
// request id has been established.
func GetRequestID(ctx context.Context) string {
	if str, ok := ctx.Value(RequestIDKey).(string); ok {
		return str
	}

	return ""
}

// SetRequestID sets the request id within a context, it will overwrite any existing
// request id.
func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}
