package templates

import (
	"log"
	"sync"
	"text/template"
)

var templates *template.Template
var once sync.Once

func Init() *template.Template {
	once.Do(func() {
		templates = template.Must(template.ParseGlob("templates/*.html"))
	})
	return templates
}

func GetTemplate() *template.Template {
	if templates == nil {
		log.Fatal("Templates not initialized. Call Init() first.")
	}
	return templates
}
