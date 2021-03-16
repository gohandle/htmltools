package htmlpage

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gohandle/htmltools/rbmux"
	"github.com/gohandle/htmltools/rbview"
	"github.com/gohandle/htmltools/rbview/rbtemplate"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// AboutHandler handles request for an about page
func AboutHandler(v rbview.V, m rbmux.M) http.HandlerFunc {
	type Page struct{ Title string }

	return func(w http.ResponseWriter, r *http.Request) {
		p := Page{"foo" + m.CurrentRoute(r)}

		if err := v.Render(w, r, p, rbview.Template("about.html")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// SetupRoutes will configure the mux router
func SetupRoutes(r *mux.Router, hf http.HandlerFunc) {
	r.Name("foo").Path("/x/y/z").Handler(hf)
}

//go:embed about.html
var tmplFS embed.FS

func Example() {
	var mr *mux.Router
	if err := fx.New(
		fx.Populate(&mr),
		fx.Supply(tmplFS),
		fx.Provide(rbtemplate.New, rbtemplate.ParseConf, rbtemplate.FromEmbed),
		fx.Provide(rbmux.New, rbmux.Helper, mux.NewRouter),
		fx.Provide(rbview.New, rbview.TemplateEncoder, AboutHandler, zap.NewDevelopment, rbview.ParseConf),
		fx.Invoke(SetupRoutes),
	).Start(context.Background()); err != nil {
		panic(err)
	}

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/x/y/z", nil)
	mr.ServeHTTP(w, r)

	fmt.Println(w.Body.String())
	//output: <h1>foofoo/x/y/z</h1>
}
