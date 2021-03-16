package rbmux

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestRoutes(t *testing.T) {
	var m M
	var mr *mux.Router
	fxtest.New(t,
		fx.Populate(&m, &mr),
		fx.Provide(New, zap.NewDevelopment, mux.NewRouter),
		fx.Invoke(func(r *mux.Router) {
			r.Name("foo").Path("/x/{y}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "%s,%s", m.CurrentRoute(r), m.URLParam(r, "y"))
			})
		}),
	).RequireStart().RequireStop()

	if act := m.URL("foo", "y", "bar").String(); act != `/x/bar` {
		t.Fatalf("got: %v", act)
	}

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/x/rab", nil)
	mr.ServeHTTP(w, r)

	if act := w.Body.String(); act != `foo,rab` {
		t.Fatalf("got: %v", act)
	}
}
