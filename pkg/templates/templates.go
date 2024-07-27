package templates

import (
	"fmt"
	Types "forum/pkg/types"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var templates *template.Template
var once sync.Once

// parseTemplates walks through the directory structure and parses all .html files
func parseTemplates(rootDir string) (*template.Template, error) {
	tmpl := template.New("")
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			_, err := tmpl.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return tmpl, err
}

// Init initializes the templates
func Init() *template.Template {
	once.Do(func() {
		var err error
		templates, err = parseTemplates("templates")
		if err != nil {
			log.Fatal("Failed to parse templates:", err)
		}
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

	err := GetTemplate().ExecuteTemplate(w, "error.html", data)
	if err != nil {
		log.Println("Failed to execute template: error.html")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LoginTemplateHandler(w http.ResponseWriter, r *http.Request) {
	err := GetTemplate().ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		log.Println("Failed to execute template: login.html")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RegisterTemplateHandler(w http.ResponseWriter, r *http.Request) {
	err := GetTemplate().ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		log.Println("Failed to execute template: register.html")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ExecuteTemplateByName(w http.ResponseWriter, name string, data interface{}) {
	err := GetTemplate().ExecuteTemplate(w, fmt.Sprintf("%s.html", name), data)
	if err != nil {
		log.Printf("Failed to execute template: %s.html", name)
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
