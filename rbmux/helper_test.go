package rbmux

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestHelpers(t *testing.T) {
	var fm template.FuncMap
	fxtest.New(t,
		fx.Populate(&fm),
		fx.Invoke(func(r *mux.Router) { r.Name("foo").Path("/x/y/z") }),
		fx.Provide(NewHelper, New, zap.NewDevelopment, mux.NewRouter)).RequireStart().RequireStop()

	tmpl := template.New("root").Funcs(fm)
	tmpl = template.Must(tmpl.Parse(`{{url "foo"}}`))

	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, nil); err != nil {
		t.Fatalf("got: %v", err)
	}

	if buf.String() != "/x/y/z" {
		t.Fatalf("got: %v", buf.String())
	}
}
