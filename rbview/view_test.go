package rbview

import (
	"embed"
	"net/http/httptest"
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"

	"github.com/gohandle/htmltools/rbview/rbtemplate"
)

var _ Encoder = &JSON{}
var _ Encoder = &Tmpl{}

func TestNoEncoders(t *testing.T) {
	var v *View
	fxtest.New(t,
		fx.Populate(&v),
		fx.Provide(New, zap.NewDevelopment, ParseConf),
	).RequireStart().RequireStop()

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	if err := v.Render(w, r, "foo"); err == nil || err.Error() != "no encoders configured" {
		t.Fatalf("got: %v", err)
	}
}

//go:embed test.html
var testFS embed.FS

func TestTemplateJSON(t *testing.T) {
	var v *View
	fxtest.New(t,
		fx.Populate(&v),
		fx.Provide(rbtemplate.New, rbtemplate.FromEmbed, rbtemplate.ParseConf), fx.Supply(testFS),
		fx.Provide(New, zap.NewDevelopment, ParseConf, JSONEncoder, TemplateEncoder),
	).RequireStart().RequireStop()

	t.Run("default encoding", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		if err := v.Render(w, r, "foo"); err != nil {
			t.Fatalf("got: %v", err)
		}

		if act := w.Body.String(); act != `"foo"`+"\n" {
			t.Fatalf("got: %v", act)
		}
	})

	t.Run("matched encoding", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		r.Header.Set("Accept", "text/html")
		if err := v.Render(w, r, "foo", Template("test.html")); err != nil {
			t.Fatalf("got: %v", err)
		}

		if act := w.Body.String(); act != `<p>foo</p>` {
			t.Fatalf("got: %v", act)
		}
	})
}
