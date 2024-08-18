package templates

import (
	"forum/pkg/handlers/categories"
	"forum/pkg/handlers/posts"
	Types "forum/pkg/types"
	"forum/pkg/utils"
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

func IndexTemplateHandler(w http.ResponseWriter, r *http.Request) {
	categories := categories.GetCategoriesWithSubcategoriesHandler(w, r)
	data := map[string]interface{}{
		"Categories": categories,
	}

	err := GetTemplate().ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Println("Failed to execute template: index.html")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SubcategoryTemplateHandler(w http.ResponseWriter, r *http.Request) {
	posts := posts.GetPostsFromSubCategoryHandler(w, r)

	isAuthenticated := utils.IsAuthenticated(r)

	data := map[string]interface{}{
		"Posts":           posts,
		"Subcategory":     posts[0].Subcategory,
		"Category":        posts[0].Subcategory.Category,
		"IsAuthenticated": isAuthenticated,
	}

	err := GetTemplate().ExecuteTemplate(w, "subcategory.html", data)
	if err != nil {
		log.Println("Failed to execute template: subcategory.html")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewPostTemplateHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Replace this with a proper handler that only gets the subcategory and category
	posts := posts.GetPostsFromSubCategoryHandler(w, r)

	data := map[string]interface{}{
		"Subcategory": posts[0].Subcategory,
		"Category":    posts[0].Subcategory.Category,
	}

	err := GetTemplate().ExecuteTemplate(w, "new-post.html", data)
	if err != nil {
		log.Println("Failed to execute template: new-post.html")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func PostTemplateHandler(w http.ResponseWriter, r *http.Request) {
	post := posts.GetPostHandler(w, r)

	isAuthenticated := utils.IsAuthenticated(r)

	var user Types.User

	if isAuthenticated {
		user = utils.GetUserInfoBySession(w, r)
	}

	data := map[string]interface{}{
		"Post":            post,
		"Subcategory":     post.Subcategory,
		"Category":        post.Subcategory.Category,
		"Comments":        post.Comments,
		"IsAuthenticated": isAuthenticated,
		"User":            user,
	}

	err := GetTemplate().ExecuteTemplate(w, "post.html", data)
	if err != nil {
		log.Println("Failed to execute template: post.html")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
