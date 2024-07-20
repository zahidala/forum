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
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	hashedPassword, hashErr := utils.HashPassword(password)
	if hashErr != nil {
		log.Println(hashErr)
		http.Error(w, "Error generating password hash", http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO users (name, username, email, password) VALUES (?, ?, ?, ?)"

	userAddExecErr := db.PrepareAndExecute(query,
		name,
		username,
		email,
		hashedPassword,
	)
	if userAddExecErr != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	fmt.Println(username, email, hashedPassword)
}

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	userQuery := "SELECT id, password FROM users WHERE username = ?"

	userStmt, userErr := db.GetDB().Prepare(userQuery)
	if userErr != nil {
		http.Error(w, "Error preparing query", http.StatusInternalServerError)
		return
	}
	defer userStmt.Close()

	var userId int
	var hashedPassword string

	findUserErr := userStmt.QueryRow(username).Scan(&userId, &hashedPassword)
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

	createdAt := time.Now()
	expiresAt := createdAt.Add(24 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:    "sessionId",
		Value:   sessionId,
		Expires: expiresAt,
	})

	sessionsAddQuery := "INSERT INTO sessions (id, userId, createdAt, expiresAt) VALUES (?, ?, ?, ?)"

	sessionAddExecErr := db.PrepareAndExecute(sessionsAddQuery, sessionId, userId, createdAt, expiresAt)
	if sessionAddExecErr != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}
}
