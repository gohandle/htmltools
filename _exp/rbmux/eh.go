package rbmux

import "net/http"

// HandlerFunc is like http.HandlerFunc but allows handlers to error
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// ServeHTTP will allow the HandlerFunc to be used as any http.Handler
func (hf HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = hf(w, r)
}
