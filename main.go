package main

import (
	"forum/pkg/db"
	"forum/pkg/handlers/users"
	"forum/pkg/middlewares"
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

	http.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		err := templates.GetTemplate().ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			log.Println(err)
			return
		}
	})

	http.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		users.UserLoginHandler(w, r)
	})

	http.HandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {
		err := templates.GetTemplate().ExecuteTemplate(w, "register.html", nil)
		if err != nil {
			log.Println(err)
			return
		}
	})

	http.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		users.CreateUserHandler(w, r)
	})

	// http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
	// 	err := templates.GetTemplate().ExecuteTemplate(w, "index.html", nil)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// })

	// An example of using the AuthRequired middleware to protect the index page

	http.Handle("GET /", middlewares.AuthRequired(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := templates.GetTemplate().ExecuteTemplate(w, "index.html", nil)

		if err != nil {
			log.Println(err)
			return
		}
	})))

	log.Println("Connected to the database")
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
