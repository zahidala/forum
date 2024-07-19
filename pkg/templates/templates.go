package templates

import "text/template"

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

func GetTemplate() *template.Template {
	return templates
}
