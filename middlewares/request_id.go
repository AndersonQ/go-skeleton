package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type reqIDCtx int

const reqIDContexKey reqIDCtx = reqIDCtx(0)

// requestIDContext creates a context with request id
func requestIDContext(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, reqIDContexKey, rid)
}

// requestIDFromContext returns the request id from context
func requestIDFromContext(ctx context.Context) string {
	rid, ok := ctx.Value(reqIDContexKey).(string)
	if !ok {
		return ""
	}
	return rid
}

// RequestIDHandler sets unique request id.
// If header `X-Request-ID` is already present in the request, that is considered the
// request id. Otherwise, generates a new unique ID.
func RequestIDHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get("X-Request-ID")
		if rid == "" {
			rid = uuid.New().String()
			r.Header.Set("X-Request-ID", rid)
		}
		ctx := requestIDContext(r.Context(), rid)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
