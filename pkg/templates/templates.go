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

func MergeBaseData(w http.ResponseWriter, r *http.Request, specificData map[string]interface{}) map[string]interface{} {
	var user Types.User

	isAuthenticated := utils.IsAuthenticated(r)

	if isAuthenticated {
		user = utils.GetUserInfoBySession(w, r)
	}

	baseData := map[string]interface{}{
		"NavbarOptions":   categories.GetCategoriesHandler(w, r),
		"IsAuthenticated": isAuthenticated,
		"User":            user,
	}

	for key, value := range specificData {
		baseData[key] = value
	}

	return baseData
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	var data Types.Error
	LoginTemplateHandler(w, r, data)
}

func LoginTemplateHandler(w http.ResponseWriter, r *http.Request, data Types.Error) {
	dataMap := map[string]interface{}{
		"Obj":   data,
		"Title": "Login",
	}

	err := GetTemplate().ExecuteTemplate(w, "login.html", MergeBaseData(w, r, dataMap))
	if err != nil {
		log.Println("Failed to execute template: login.html")
		log.Println(err)

		ErrorTemplateHandler(w, r, Types.ErrorPageProps{
			Error: Types.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			},
			Title: "Internal Server Error",
		})
		return
	}
}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	var data Types.RegValidation
	RegisterTemplateHandler(w, r, data)
}

func RegisterTemplateHandler(w http.ResponseWriter, r *http.Request, data Types.RegValidation) {
	dataMap := map[string]interface{}{
		"Obj":   data,
		"Title": "Register",
	}

	err := GetTemplate().ExecuteTemplate(w, "register.html", MergeBaseData(w, r, dataMap))
	if err != nil {
		log.Println("Failed to execute template: register.html")
		log.Println(err)

		ErrorTemplateHandler(w, r, Types.ErrorPageProps{
			Error: Types.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			},
			Title: "Internal Server Error",
		})
		return
	}
}

func IndexTemplateHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorTemplateHandler(w, r, Types.ErrorPageProps{
			Error: Types.Error{
				Message: "Page not found.",
				Code:    404,
			},
			Title: "Page Not Found",
		})
		return
	}

	categories := categories.GetCategoriesHandler(w, r)
	newPosts := posts.GetNewPostsHandler(w, r)
	allPosts := posts.GetAllPostsHandler(w, r)
	isAuthenticated := utils.IsAuthenticated(r)

	data := map[string]interface{}{
		"Categories":      categories,
		"NewPosts":        newPosts,
		"AllPosts":        allPosts,
		"IsAuthenticated": isAuthenticated,
	}

	err := GetTemplate().ExecuteTemplate(w, "index.html", MergeBaseData(w, r, data))
	if err != nil {
		log.Println("Failed to execute template: index.html")
		log.Println(err)

		ErrorTemplateHandler(w, r, Types.ErrorPageProps{
			Error: Types.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			},
			Title: "Internal Server Error",
		})
		return
	}
}

func CategoryTemplateHandler(w http.ResponseWriter, r *http.Request) {
	category, error := categories.GetCategoryHandler(w, r)

	if error.Error.Code != 200 && error.Error.Code != 0 {
		ErrorTemplateHandler(w, r, error)
		return
	}

	posts := posts.GetPostsFromCategoryHandler(w, r)

	data := map[string]interface{}{
		"Posts":    posts,
		"Category": category,
		"Title":    category.Name,
	}

	err := GetTemplate().ExecuteTemplate(w, "category.html", MergeBaseData(w, r, data))
	if err != nil {
		log.Println("Failed to execute template: category.html")
		log.Println(err)

		ErrorTemplateHandler(w, r, Types.ErrorPageProps{
			Error: Types.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			},
			Title: "Internal Server Error",
		})
		return
	}
}

func NewPostTemplateHandler(w http.ResponseWriter, r *http.Request) {

	categories := categories.GetCategoriesHandler(w, r)

	data := map[string]interface{}{
		"Categories": categories,
		"Title":      "New Post",
	}

	err := GetTemplate().ExecuteTemplate(w, "new-post.html", MergeBaseData(w, r, data))
	if err != nil {
		log.Println("Failed to execute template: new-post.html")
		log.Println(err)

		ErrorTemplateHandler(w, r, Types.ErrorPageProps{
			Error: Types.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			},
			Title: "Internal Server Error",
		})
		return
	}
}

func NewPostByCategoryTemplateHandler(w http.ResponseWriter, r *http.Request) {
	category, error := categories.GetCategoryHandler(w, r)

	if error.Error.Code != 200 && error.Error.Code != 0 {
		ErrorTemplateHandler(w, r, error)
		return
	}

	data := map[string]interface{}{
		"Category": category,
		"Title":    "New Post - " + category.Name,
	}

	err := GetTemplate().ExecuteTemplate(w, "new-post-by-category.html", data)
	if err != nil {
		log.Println("Failed to execute template: new-post-by-category.html")
		log.Println(err)

		ErrorTemplateHandler(w, r, Types.ErrorPageProps{
			Error: Types.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			},
			Title: "Internal Server Error",
		})
		return
	}
}

func PostTemplateHandler(w http.ResponseWriter, r *http.Request) {
	post, postErr := posts.GetPostHandler(w, r)

	if postErr.Error.Code != 200 && postErr.Error.Code != 0 {
		ErrorTemplateHandler(w, r, postErr)
		return
	}

	postLikes := posts.GetPostLikesHandler(w, r)
	postDislikes := posts.GetPostDislikesHandler(w, r)
	comments := comments.GetCommentsHandler(w, r)

	categories := categories.GetCategoriesPostHandler(w, r, post.ID)

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
		"Categories":                  categories,
		"Comments":                    comments,
		"IsPostLikedByCurrentUser":    isPostLikedByCurrentUser,
		"IsPostDislikedByCurrentUser": isPostDislikedByCurrentUser,
		"Title":                       post.Title,
	}

	err := GetTemplate().ExecuteTemplate(w, "post.html", MergeBaseData(w, r, data))
	if err != nil {
		log.Println("Failed to execute template: post.html")
		log.Println(err)

		ErrorTemplateHandler(w, r, Types.ErrorPageProps{
			Error: Types.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			},
			Title: "Internal Server Error",
		})
		return
	}

}

// ErrorTemplateHandler renders the error page
func ErrorTemplateHandler(w http.ResponseWriter, r *http.Request, data Types.ErrorPageProps) {
	w.WriteHeader(data.Error.Code)

	err := GetTemplate().ExecuteTemplate(w, "error.html", data)
	if err != nil {
		log.Println("Failed to execute template: error.html")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
