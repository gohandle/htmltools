package rblog

import (
	"context"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// make context keys private
type ctxKey string

// Logging annotates the logging middleware as rb middlware
var Logging = fx.Annotated{
	Target: NewLogging,
	Name:   "rb.ware.logging",
}

// NewLogging will log requests and configure a request scoped logger for
// the request
func NewLogging(logs *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := logs.With(
				zap.String("request_url", r.URL.String()),
				zap.String("request_method", r.Method),
				zap.Int64("request_length", r.ContentLength),
				zap.String("request_host", r.Host),
				zap.String("request_uri", r.RequestURI),
			)

			next.ServeHTTP(w, r.WithContext(
				WithLogger(r.Context(), l)))
		})
	}
}

// WithLogger sets request scoped logger
func WithLogger(ctx context.Context, logs *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey("logger"), logs)
}

// Logger returns the request scoped logger, returns nil if none is configured
func Logger(ctx context.Context) (l *zap.Logger) {
	l, _ = ctx.Value(ctxKey("logger")).(*zap.Logger)
	return
}

// L is a convenience method for retrieving a logger from the request's context. If non
// is configured it will return a nop logger.
func L(r *http.Request) (l *zap.Logger) {
	if l = Logger(r.Context()); l != nil {
		return l
	}

	return zap.NewNop()
}
