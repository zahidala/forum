package main

import (
	"forum/pkg/db"
	"forum/pkg/handlers/users"
	"forum/pkg/templates"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
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

	http.HandleFunc("GET /post/{id}", templates.PostTemplateHandler)

	// An example of using the AuthRequired middleware to protect the index page

	// http.Handle("GET /", middlewares.AuthRequired(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	templates.ExecuteTemplateByName(w, "index", nil)
	// })))

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
