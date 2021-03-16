package helloworld

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gohandle/htmltools/rbbind"
	"github.com/gohandle/htmltools/rbview"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func EchoHandler(b *rbbind.Binder, v *rbview.View) http.HandlerFunc {
	type Input struct{ Name string }
	type Output struct{ Message string }

	return func(w http.ResponseWriter, r *http.Request) {
		var in Input
		if err := b.Bind(r, &in); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		out := Output{"Hello, " + strings.ToUpper(in.Name)}
		if err := v.Render(w, r, out); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func Example() {
	ctx := context.Background()

	// we use depdency injection to build our handler
	var h http.HandlerFunc
	app := fx.New(fx.Provide(
		EchoHandler, rbview.New, rbbind.New, zap.NewDevelopment,
		rbview.ParseConf, rbbind.ParseConf, rbbind.JSONDecoder, rbview.JSONEncoder,
	), fx.Populate(&h))

	// we start the application to resolve all dependencies
	if err := app.Start(ctx); err != nil {
		panic(err)
	}

	// create a request and let our handler handle it
	b := strings.NewReader(`{"name":"foo"}`)
	w, r := httptest.NewRecorder(), httptest.NewRequest("POST", "/foo", b)
	r.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(w, r)

	fmt.Println(w.Body.String())

	// shut down our application
	if err := app.Stop(ctx); err != nil {
		panic(err)
	}

	//output: hello, FOO
}
