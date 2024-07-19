package main

import (
	"forum/pkg/db"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize the database connection
	db.Init()
	defer db.CloseDB()

	log.Println("Connected to the database")
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
