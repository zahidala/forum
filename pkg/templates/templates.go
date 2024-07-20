package templates

import (
	Types "forum/pkg/types"
	"log"
	"net/http"
	"sync"
	"text/template"
)

var templates *template.Template
var once sync.Once

// Init initializes the templates
func Init() *template.Template {
	once.Do(func() {
		templates = template.Must(template.ParseGlob("templates/*.html"))
	})
	return templates
}

// GetTemplate returns the template instance
func GetTemplate() *template.Template {
	if templates == nil {
		log.Fatal("Templates not initialized. Call Init() first.")
	}
	return templates
}

// ErrorTemplate renders the error page
func ErrorTemplate(w http.ResponseWriter, data *Types.ErrorPageProps) {
	w.WriteHeader(data.Error.Code)

	err := GetTemplate().ExecuteTemplate(w, "error.html", nil)
	if err != nil {
		log.Println(err)
		return
	}
}
