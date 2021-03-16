package rbbind

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

var _ Decoder = &Form{}

func TestBinder(t *testing.T) {
	var b B
	fxtest.New(t,
		fx.Populate(&b),
		fx.Provide(ParseConf, FormDecoder, JSONDecoder, New, zap.NewDevelopment),
	).RequireStart().RequireStop()

	t.Run("decode form query", func(t *testing.T) {
		var v struct {
			Foo string `form:"foo"`
		}

		r := httptest.NewRequest("GET", "/?foo=bar", nil)
		if err := b.Bind(r, &v); err != nil || v.Foo != "bar" {
			t.Fatalf("got: '%v' '%v'", err, v.Foo)
		}
	})

	t.Run("decode form body", func(t *testing.T) {
		var v struct {
			Foo string `form:"foo"`
			Bar string `form:"bar"`
		}

		r := httptest.NewRequest("POST", "/?foo=bar", strings.NewReader(`bar=rab`))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err := b.Bind(r, &v); err != nil || v.Foo != "bar" || v.Bar != "rab" {
			t.Fatalf("got: '%v' '%v' '%v'", err, v.Foo, v.Bar)
		}
	})

	t.Run("decode form query, json body", func(t *testing.T) {
		var v struct {
			Foo string `form:"foo"`
			Bar string `form:"bar" json:"bar"`
		}

		r := httptest.NewRequest("POST", "/?foo=xyz&bar=def", strings.NewReader(`{"bar": "ddd"}`))
		r.Header.Set("content-type", "application/json")

		if err := b.Bind(r, &v); err != nil || v.Foo != "xyz" || v.Bar != "ddd" {
			t.Fatalf("got: '%v' '%v' '%v'", err, v.Foo, v.Bar)
		}
	})
}

func TestBinderWithOrder(t *testing.T) {
	os.Setenv("RB_BIND_DECODER_ORDER", "json,form")
	defer os.Unsetenv("RB_BIND_DECODER_ORDER")

	var b B
	fxtest.New(t,
		fx.Populate(&b),
		fx.Provide(ParseConf, FormDecoder, JSONDecoder, New, zap.NewDevelopment),
	).RequireStart().RequireStop()

	var v struct {
		Foo string `form:"foo"`
		Bar string `form:"bar" json:"bar"`
	}

	r := httptest.NewRequest("POST", "/?foo=xyz&bar=def", strings.NewReader(`{"bar": "ddd"}`))
	r.Header.Set("content-type", "application/json")

	if err := b.Bind(r, &v); err != nil || v.Foo != "xyz" || v.Bar != "def" {
		t.Fatalf("got: '%v' '%v' '%v'", err, v.Foo, v.Bar)
	}
}
