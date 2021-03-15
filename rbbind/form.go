package rbbind

import (
	"fmt"
	"net/http"

	"github.com/go-playground/form/v4"
	"go.uber.org/fx"
)

type Form struct{ dec *form.Decoder }

func NewForm() Decoder {
	return &Form{dec: form.NewDecoder()}
}

func (d *Form) Name() string { return "form" }
func (d *Form) Decode(r *http.Request, v interface{}, _ string, _ map[string]string) (err error) {
	if err = r.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form: %w", err)
	}

	return d.dec.Decode(v, r.Form)
}

var FormDecoder = fx.Annotated{
	Target: NewForm,
	Group:  "rb.decoder",
}
