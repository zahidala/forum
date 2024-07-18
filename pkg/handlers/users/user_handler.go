package users

import (
	"database/sql"
	"fmt"
	"forum/pkg/db"
	"forum/pkg/utils"
	"log"
	"net/http"
	"time"
)

func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation specific to fetching users

	// cookie, cookieErr := r.Cookie("session_id")

	// if cookieErr == nil {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }

	// sessionId := cookie.Value
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	hashedPassword, hashErr := utils.HashPassword(password)
	if hashErr != nil {
		log.Println(hashErr)
		http.Error(w, "Error generating password hash", http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error preparing query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, addUserErr := stmt.Exec(username, email, hashedPassword)
	if addUserErr != nil {
		log.Println(err)
		http.Error(w, "Error executing query", http.StatusInternalServerError)
		return
	}

	fmt.Println(username, email, hashedPassword)
}

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	userQuery := "SELECT password FROM users WHERE username = ?"

	userStmt, userErr := db.GetDB().Prepare(userQuery)
	if userErr != nil {
		http.Error(w, "Error preparing query", http.StatusInternalServerError)
		return
	}
	defer userStmt.Close()

	var hashedPassword string

	findUserErr := userStmt.QueryRow(username).Scan(&hashedPassword)
	switch {
	case findUserErr == sql.ErrNoRows:
		http.Error(w, "User not found", http.StatusNotFound)
		return
	case findUserErr != nil:
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}

	compareErr := utils.CompareHashAndPassword(hashedPassword, password)
	if compareErr != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	sessionId := utils.GenerateSessionID()

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   sessionId,
		Expires: time.Now().Add(24 * time.Hour),
	})

	sessionsQuery := "INSERT INTO sessions (id, userId, createdAt, expiresAt) VALUES (?, ?, ?, ?)"

	sessionStmt, sessionErr := db.GetDB().Prepare(sessionsQuery)
	if sessionErr != nil {
		http.Error(w, "Error preparing query", http.StatusInternalServerError)
		return
	}
	defer sessionStmt.Close()
}
