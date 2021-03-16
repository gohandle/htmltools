package rbview

import (
	"fmt"
	"html/template"
	"net/http"

	"go.uber.org/fx"
)

type Tmpl struct{ tmpl *template.Template }

func NewTemplate(tmpl *template.Template) Encoder { return &Tmpl{tmpl} }

func (e *Tmpl) MIME() string { return "text/html" }
func (e *Tmpl) Encode(w http.ResponseWriter, r *http.Request, v interface{}, o Options) error {
	if o.TemplateName == "" {
		return fmt.Errorf("template encoder requires Template option to be specified")
	}

	return e.tmpl.ExecuteTemplate(w, o.TemplateName, v)
}

var TemplateEncoder = fx.Annotated{
	Target: NewTemplate,
	Group:  "rb.encoder",
}
