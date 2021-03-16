package rbmux

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// M is an interface for reading mux information
type M interface {
	CurrentRoute(req *http.Request) string
	URLParam(r *http.Request, k string) string
	URL(name string, pairs ...string) *url.URL
}

// Mux implements the m interface
type Mux struct{ mux *mux.Router }

// New creates a router instance
func New(logs *zap.Logger, mr *mux.Router) M {
	return &Mux{mr}
}

// CurrentRoute returns the name of the current route
func (m *Mux) CurrentRoute(req *http.Request) (name string) {
	return mux.CurrentRoute(req).GetName()
}

// Param returns the url parameter value with name 'name', or an empty string
func (m *Mux) URLParam(req *http.Request, name string) string {
	return mux.Vars(req)[name]
}

// URL generates a url from the route with name 'name' and each parameter provides as key, value in 'pairs'
func (m *Mux) URL(name string, pairs ...string) *url.URL {
	route := m.mux.Get(name)
	if route == nil {
		// @TODO error if route doesn't exist
		return new(url.URL)
	}

	loc, err := route.URL(pairs...)
	if err != nil {
		// @TODO log error
		return new(url.URL)
	}

	return loc
}
