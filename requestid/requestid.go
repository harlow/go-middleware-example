package requestid

import (
	"context"
	"errors"
	"net/http"
)

// FromRequest request ID from http headers
func FromRequest(req *http.Request) (string, error) {
	requestID := req.Header.Get("X-Request-Id")
	if requestID == "" {
		return "", errors.New("Request ID not provided")
	}
	return requestID, nil
}

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int

// requestIDKey is the context key for the Request ID.  Its value of zero is
// arbitrary. If this package defined other context keys, they would have
// different integer values.
const requestIDKey key = 0

// NewContext returns a new Context carrying userIP.
func NewContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// FromContext extracts the user IP address from ctx, if present.
func FromContext(ctx context.Context) (string, bool) {
	// ctx.Value returns nil if ctx has no value for the key;
	// the ctx.Value type assertion returns ok=false for nil.
	requestID, ok := ctx.Value(requestIDKey).(string)
	return requestID, ok
}
