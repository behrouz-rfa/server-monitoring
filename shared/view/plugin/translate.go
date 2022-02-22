package plugin

import (
	"github.com/leonelquinteros/gotext"
	"html/template"
	"server-monitoring/shared/view"
)

//transale plugin
func Translate(v view.View) template.FuncMap {
	f := make(template.FuncMap)
	f["gettext"] = func(s string) string {
		tr := gotext.Get(s)
		return tr
	}

	return f
}
