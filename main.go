package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if _, err := os.Stat("forum.db"); os.IsNotExist(err) {
		log.Println("forum.db does not exist!")
		return
	}

	log.Println("Connected to the database")
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
