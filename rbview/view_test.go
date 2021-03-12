package rbview

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"text/template"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestParseConf(t *testing.T) {
	os.Setenv("VIEW_DIR", "foo")
	defer os.Unsetenv("VIEW_DIR")
	os.Setenv("VIEW_PATTERNS", "index.html:foo.html")
	defer os.Unsetenv("VIEW_PATTERNS")

	cfg, err := ParseConf()
	if !reflect.DeepEqual(cfg, Conf{
		Dir:      "foo",
		Patterns: []string{"index.html", "foo.html"},
	}) {
		t.Fatalf("got: %+v err: %v", cfg, err)
	}
}

func TestDirView(t *testing.T) {
	dir, _ := os.MkdirTemp("", "")
	os.WriteFile(filepath.Join(dir, "index.html"), []byte(`{{.}}`), 0777)

	var v *template.Template
	fxtest.New(t,
		fx.Supply(Conf{Dir: dir, Patterns: []string{"*.html"}}),
		fx.Provide(zap.NewDevelopment, FromDir, New),
		fx.Populate(&v)).RequireStart().RequireStop()

	buf := bytes.NewBuffer(nil)
	if err := v.ExecuteTemplate(buf, "index.html", "foo"); err != nil {
		t.Fatalf("got: %v", err)
	}

	if act := buf.String(); act != `foo` {
		t.Fatalf("got: %v", act)
	}
}
