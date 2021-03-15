package rbasset

import (
	"embed"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestAssetMiddleware(t *testing.T) {
	var sw StaticWare
	fxtest.New(t,
		fx.Populate(&sw),
		fx.Provide(NewStatic, zap.NewDevelopment, ParseConf, FromDir),
	).RequireStart().RequireStop()

	h := sw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, world")
	}))

	w1, r1 := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	h.ServeHTTP(w1, r1)
	if w1.Code != 200 {
		t.Fatalf("got: %v", w1.Code)
	}

	w2, r2 := httptest.NewRecorder(), httptest.NewRequest("GET", "/static/x.txt", nil)
	h.ServeHTTP(w2, r2)
	if w2.Code != 404 {
		t.Fatalf("got: %v", w1.Code)
	}

	w3, r3 := httptest.NewRecorder(), httptest.NewRequest("GET", "/static/go.mod", nil)
	h.ServeHTTP(w3, r3)
	if w3.Code != 200 {
		t.Fatalf("got: %v", w3.Code)
	}
}

//go:embed go.mod
var testFs embed.FS

func TestEmbedFsMiddleware(t *testing.T) {
	var sw StaticWare
	fxtest.New(t,
		fx.Supply(testFs),
		fx.Populate(&sw),
		fx.Provide(NewStatic, zap.NewDevelopment, ParseConf, FromEmbed),
	).RequireStart().RequireStop()

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/static/go.mod", nil)
	sw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, world")
	})).ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("got: %v", w.Code)
	}
}
