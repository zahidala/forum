package main

import (
	"forum/pkg/db"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := templates.GetTemplate().ExecuteTemplate(w, "index.html", nil)

		if err != nil {
			log.Println(err)
		}
	})

	log.Println("Connected to the database")
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
