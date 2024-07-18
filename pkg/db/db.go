package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var dbErrors error

func init() {
	db, dbErrors = sql.Open("sqlite3", "./forum.db")
	if dbErrors != nil {
		log.Printf("Error opening the database: %s", dbErrors)
		return
	}

	if err := db.Ping(); err != nil {
		log.Printf("Error connecting to the database: %s", err)
		return
	}
}

func GetDB() *sql.DB {
	if db == nil {
		log.Println("Database connection is not initialized")
		return nil
	}
	return db
}
