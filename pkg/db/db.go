package db

import (
	"database/sql"
	Types "forum/pkg/types"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var instance *Types.Database
var once sync.Once

// Init initializes the database connection
func Init() {
	once.Do(func() {
		conn, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
			log.Fatalf("Error opening the database: %s", err)
			return
		}

		if err := conn.Ping(); err != nil {
			log.Fatalf("Error connecting to the database: %s", err)
			return
		}

		instance = &Types.Database{
			Conn: conn,
		}
	})
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	if instance == nil {
		log.Fatal("Database not initialized. Call Init() first.")
	}
	return instance.Conn
}

// CloseDB closes the database connection
func CloseDB() {
	if instance != nil {
		instance.Mu.Lock()
		defer instance.Mu.Unlock()
		if err := instance.Conn.Close(); err != nil {
			log.Printf("Error closing the database: %s", err)
		}
	}
}

func PrepareAndExecute(query string, args ...interface{}) error {
	stmt, stmtErr := GetDB().Prepare(query)
	if stmtErr != nil {
		log.Printf("Error preparing query: %s", stmtErr)
		return stmtErr
	}

	defer stmt.Close()

	_, execErr := stmt.Exec(args...)
	if execErr != nil {
		log.Printf("Error executing query: %s", execErr)
		return execErr
	}

	return nil
}
