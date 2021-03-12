package rbhelper

import (
	"html/template"
	"time"

	"go.uber.org/fx"
)

// NewNow creates a helper
func NewNow() template.FuncMap {
	return template.FuncMap{
		"now": func() time.Time { return time.Now() },
	}
}

// NowHelper annotates as helper
var NowHelper = fx.Annotated{
	Target: NewNow,
	Group:  "rbhelper",
}
