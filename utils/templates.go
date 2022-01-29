package utils

import (
	"html/template"
	"net/http"
)

var templates *template.Template

func LoadTemplates(patters string){
	templates = template.Must(template.ParseGlob(patters))
}

func ExecuteTemplate(w http.ResponseWriter, tmpl string, data interface{}){
	templates.ExecuteTemplate(w,tmpl,data)
}
