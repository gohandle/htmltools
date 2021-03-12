package rbware

import (
	"net/http"

	"go.uber.org/zap"
)

// Logging logs all requests and provides a request scoped zap.Logger
func Logging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

// L
func L(r *http.Request) *zap.Logger {
	return zap.NewNop()
}
