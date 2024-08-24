package templates

import (
	"forum/pkg/handlers/categories"
	"forum/pkg/handlers/comments"
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
func ErrorTemplate(w http.ResponseWriter, data Types.ErrorPageProps) {
	w.WriteHeader(data.Error.Code)

	err := GetTemplate().ExecuteTemplate(w, "error.html", data)
	if err != nil {
		log.Println("Failed to execute template: error.html")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	var data Types.Error
	LoginTemplateHandler(w, r, data)
}

func LoginTemplateHandler(w http.ResponseWriter, r *http.Request, data Types.Error) {
	err := GetTemplate().ExecuteTemplate(w, "login.html", data)
	if err != nil {
		log.Println("Failed to execute template: login.html")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	var data Types.RegValidation
	RegisterTemplateHandler(w, r, data)
}

func RegisterTemplateHandler(w http.ResponseWriter, r *http.Request, data Types.RegValidation) {
	err := GetTemplate().ExecuteTemplate(w, "register.html", data)
	if err != nil {
		log.Println("Failed to execute template: register.html")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func IndexTemplateHandler(w http.ResponseWriter, r *http.Request) {
	categories := categories.GetCategoriesHandler(w, r)
	newPosts := posts.GetNewPostsHandler(w, r)

	data := map[string]interface{}{
		"Categories": categories,
		"NewPosts":   newPosts,
	}

	err := GetTemplate().ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Println("Failed to execute template: index.html")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CategoryTemplateHandler(w http.ResponseWriter, r *http.Request) {
	category := categories.GetCategoryHandler(w, r)

	posts := posts.GetPostsFromCategoryHandler(w, r)

	isAuthenticated := utils.IsAuthenticated(r)

	data := map[string]interface{}{
		"Posts":           posts,
		"Category":        category,
		"IsAuthenticated": isAuthenticated,
	}

	err := GetTemplate().ExecuteTemplate(w, "category.html", data)
	if err != nil {
		log.Println("Failed to execute template: category.html")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		ErrorTemplateHandler(w, r, Types.ErrorPageProps{
			Error: Types.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			},
		})
		return
	}
}

func NewPostTemplateHandler(w http.ResponseWriter, r *http.Request) {
	category := categories.GetCategoryHandler(w, r)

	var user Types.User

	isAuthenticated := utils.IsAuthenticated(r)

	if isAuthenticated {
		user = utils.GetUserInfoBySession(w, r)
	}

	data := map[string]interface{}{
		"Category": category,
		"User":     user,
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
	postLikes := posts.GetPostLikesHandler(w, r)
	postDislikes := posts.GetPostDislikesHandler(w, r)
	comments := comments.GetCommentsHandler(w, r)

	category := categories.GetCategoryHandler(w, r)

	isAuthenticated := utils.IsAuthenticated(r)

	var user Types.User
	isPostLikedByCurrentUser := false
	isPostDislikedByCurrentUser := false

	if isAuthenticated {
		user = utils.GetUserInfoBySession(w, r)
		isPostLikedByCurrentUser = posts.IsPostLikedByCurrentUserHandler(w, r, post.ID, user.ID)
		isPostDislikedByCurrentUser = posts.IsPostDisLikedByCurrentUserHandler(w, r, post.ID, user.ID)
	}

	data := map[string]interface{}{
		"Post":                        post,
		"Likes":                       postLikes,
		"Dislikes":                    postDislikes,
		"Category":                    category,
		"Comments":                    comments,
		"IsAuthenticated":             isAuthenticated,
		"User":                        user,
		"IsPostLikedByCurrentUser":    isPostLikedByCurrentUser,
		"IsPostDislikedByCurrentUser": isPostDislikedByCurrentUser,
	}

	err := GetTemplate().ExecuteTemplate(w, "post.html", data)
	if err != nil {
		log.Println("Failed to execute template: post.html")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func ErrorTemplateHandler(w http.ResponseWriter, r *http.Request, data Types.ErrorPageProps) {
	err := GetTemplate().ExecuteTemplate(w, "error.html", data)
	if err != nil {
		log.Println("Failed to execute template: error.html")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
