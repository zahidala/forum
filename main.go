package main

import (
	"encoding/json"
	"forum/pkg/db"
	"forum/pkg/env"
	"forum/pkg/handlers/comments"
	"forum/pkg/handlers/posts"
	"forum/pkg/handlers/uploads"
	"forum/pkg/handlers/users"
	"forum/pkg/templates"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize the environment variables
	env.Init()

	// Initialize the database connection
	db.Init()
	defer db.CloseDB()

	// Initialize the templates
	templates.Init()

	http.Handle("GET /static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	http.HandleFunc("GET /login", templates.LoginTemplateHandler)
	http.HandleFunc("POST /login", users.UserLoginHandler)

	http.HandleFunc("GET /logout", users.UserLogoutHandler)

	http.HandleFunc("GET /register", templates.RegisterTemplateHandler)
	http.HandleFunc("POST /register", users.CreateUserHandler)

	http.HandleFunc("GET /", templates.IndexTemplateHandler)

	http.HandleFunc("GET /subcategory/{id}", templates.SubcategoryTemplateHandler)
	http.HandleFunc("GET /subcategory/{id}/new-post", templates.NewPostTemplateHandler)
	http.HandleFunc("POST /subcategory/{id}/new-post", posts.CreatePostHandler)

	http.HandleFunc("GET /post/{id}", templates.PostTemplateHandler)
	http.HandleFunc("POST /post/{id}/comment", comments.CreateCommentHandler)
	http.HandleFunc("PUT /post/{id}/like", posts.PostLikeHandler)
	http.HandleFunc("PUT /post/{id}/remove-like", posts.PostRemoveLikeHandler)
	http.HandleFunc("PUT /post/{id}/dislike", posts.PostDislikeHandler)
	http.HandleFunc("PUT /post/{id}/remove-dislike", posts.PostRemoveDislikeHandler)

	http.HandleFunc("PUT /comment/{id}/like", comments.CommentLikeHandler)
	http.HandleFunc("PUT /comment/{id}/dislike", comments.CommentDislikeHandler)
	http.HandleFunc("PUT /comment/{id}/remove-like", comments.CommentRemoveLikeHandler)
	http.HandleFunc("PUT /comment/{id}/remove-dislike", comments.CommentRemoveDislikeHandler)

	http.HandleFunc("POST /upload", func(w http.ResponseWriter, r *http.Request) {
		uploadedImage := uploads.UploadImageHandler(w, r)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(uploadedImage)
	})

	// An example of using the AuthRequired middleware to protect the index page

	// http.Handle("GET /", middlewares.AuthRequired(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	templates.ExecuteTemplateByName(w, "index", nil)
	// })))

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
