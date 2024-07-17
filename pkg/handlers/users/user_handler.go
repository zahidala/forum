package users

import (
	"database/sql"
	"fmt"
	"forum/pkg/db"
	"log"
	"net/http"

	bcrypt "golang.org/x/crypto/bcrypt"
)

func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation specific to fetching users
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	bytePass := []byte(password)

	hash, err := bcrypt.GenerateFromPassword(bytePass, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error generating password hash", http.StatusInternalServerError)
		return
	}
	hashedPassword := string(hash)

	query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error preparing query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, email, hashedPassword)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error executing query", http.StatusInternalServerError)
		return
	}

	fmt.Println(username, email, hashedPassword)
}

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	bytePass := []byte(password)

	query := "SELECT password FROM users WHERE username = ?"

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Error preparing query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var hashedPassword string

	err = stmt.QueryRow(username).Scan(&hashedPassword)
	switch {
	case err == sql.ErrNoRows:
		http.Error(w, "User not found", http.StatusNotFound)
		return
	case err != nil:
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), bytePass)
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// handle successful login
}
