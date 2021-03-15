package rblog

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"go.uber.org/zap"
)

// make context keys private
type ctxKey string

// IDHeaders hold the names of common request id header
type IDHeaders []string

// CommonIDHeaders return common id header names
func CommonIDHeaders() IDHeaders {
	return []string{
		"X-Request-ID", "X-Correlation-ID", // unofficial standards: https://en.wikipedia.org/wiki/List_of_HTTP_header_fields
		"X-Amzn-Trace-Id",       // AWS xray tracing: https://docs.aws.amazon.com/xray/latest/devguide/xray-concepts.html#xray-concepts-tracingheader
		"Cf-Request-Id",         // Cloudflare: https://community.cloudflare.com/t/new-http-response-header-cf-request-id/165869
		"X-Cloud-Trace-Context", // Google Cloud https://cloud.google.com/appengine/docs/standard/go/reference/request-response-headers
	}
}

// RandRead is used for request id generation. It can be ovewritten in test to make them fully
// deterministic. The default is set to a non-cryptographic random number
var RandRead = rand.Read

// LoggingWare is the request logger middleware
type LoggingWare func(http.Handler) http.Handler

// NewLogging will provide a request-scoped logger with a request id field.
func NewLogging(logs *zap.Logger, hdrs IDHeaders) LoggingWare {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var rid string
			for _, hdrn := range hdrs {
				rid = r.Header.Get(hdrn)
				if rid != "" {
					break
				}
			}
			if rid == "" {
				var b [18]byte
				if _, err := RandRead(b[:]); err != nil {
					logs.Error("failed to read random bytes for request id middleware",
						zap.Error(err))
				}
				rid = base64.URLEncoding.EncodeToString(b[:])
			}

			l := logs.With(zap.String("request_id", rid))
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
// is configured it will return the global zap logger
func L(r *http.Request) (l *zap.Logger) {
	if l = Logger(r.Context()); l != nil {
		return l
	}

	return zap.L()
}
