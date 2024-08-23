package main

import (
	"encoding/json"
	"forum/pkg/db"
	"forum/pkg/env"
	"forum/pkg/handlers/comments"
	"forum/pkg/handlers/posts"
	"forum/pkg/handlers/uploads"
	"forum/pkg/handlers/users"
	"forum/pkg/middlewares"
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

	http.HandleFunc("GET /login", templates.LoginPageHandler)
	http.HandleFunc("POST /login", users.UserLoginHandler)

	http.HandleFunc("GET /logout", users.UserLogoutHandler)

	http.HandleFunc("GET /register", templates.RegisterPageHandler)
	http.HandleFunc("POST /register", users.CreateUserHandler)

	http.HandleFunc("GET /", templates.IndexTemplateHandler)

	http.HandleFunc("GET /category/{id}", templates.CategoryTemplateHandler)
	http.HandleFunc("GET /category/{id}/new-post", templates.NewPostTemplateHandler)

	http.Handle("POST /category/{id}/new-post", middlewares.AuthRequired(http.HandlerFunc(posts.CreatePostHandler)))

	http.HandleFunc("GET /post/{id}", templates.PostTemplateHandler)

	http.Handle("POST /post/{id}/comment", middlewares.AuthRequired(http.HandlerFunc(comments.CreateCommentHandler)))

	http.Handle("PUT /post/{id}/like", middlewares.AuthRequired(http.HandlerFunc(posts.PostLikeHandler)))
	http.Handle("PUT /post/{id}/remove-like", middlewares.AuthRequired(http.HandlerFunc(posts.PostRemoveLikeHandler)))
	http.Handle("PUT /post/{id}/dislike", middlewares.AuthRequired(http.HandlerFunc(posts.PostDislikeHandler)))
	http.Handle("PUT /post/{id}/remove-dislike", middlewares.AuthRequired(http.HandlerFunc(posts.PostRemoveDislikeHandler)))

	http.Handle("PUT /comment/{id}/like", middlewares.AuthRequired(http.HandlerFunc(comments.CommentLikeHandler)))
	http.Handle("PUT /comment/{id}/dislike", middlewares.AuthRequired(http.HandlerFunc(comments.CommentDislikeHandler)))
	http.Handle("PUT /comment/{id}/remove-like", middlewares.AuthRequired(http.HandlerFunc(comments.CommentRemoveLikeHandler)))
	http.Handle("PUT /comment/{id}/remove-dislike", middlewares.AuthRequired(http.HandlerFunc(comments.CommentRemoveDislikeHandler)))

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
