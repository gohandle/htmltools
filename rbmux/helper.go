package rbmux

import (
	"html/template"
	"net/url"

	"go.uber.org/fx"
)

// NewHelper creates a helper
func NewHelper(m M) template.FuncMap {
	return template.FuncMap{
		"url": func(name string, pairs ...string) (*url.URL, error) {
			// @TODO report error here, if no route exists
			return m.URL(name, pairs...), nil
		},
	}
}

// Helper annotates all mux helpers as helper
var Helper = fx.Annotated{
	Target: NewHelper,
	Group:  "rb.helper",
}
