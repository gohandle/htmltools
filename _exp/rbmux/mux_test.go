package rbmux_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gohandle/htmltools/rbmux"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func AuthMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

func HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func AboutHandler() rbmux.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		return
	}
}

type Handlers struct {
	fx.In
	About rbmux.HandlerFunc `name:"handler.about"`
	Home  http.HandlerFunc  `name:"handler.home"`
}

type Middlewares struct {
	fx.In
	Auth mux.MiddlewareFunc `name:"middleware.auth"`
}

func New(r *mux.Router, h Handlers, mw Middlewares) http.Handler {
	r.Use(mw.Auth)
	r.Name("about").Path("/about").Methods("GET", "POST").Handler(h.About)
	r.Name("home").Path("/").Methods("GET", "POST").Handler(h.Home)
	return r
}

func TestMux(t *testing.T) {
	var h http.Handler
	fxtest.New(t,
		fx.Provide(zap.NewDevelopment, mux.NewRouter, New,
			fx.Annotated{Target: AuthMiddleware, Name: "middleware.auth"},
			fx.Annotated{Target: AboutHandler, Name: "handler.about"},
			fx.Annotated{Target: HomeHandler, Name: "handler.home"}),
		fx.Populate(&h)).RequireStart().RequireStop()

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/about", nil)
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("got: %v", w.Code)
	}

}
