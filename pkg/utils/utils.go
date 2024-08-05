package utils

import (
	"forum/pkg/db"
	"log"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

// GenerateSessionID generates a new UUID for session ID
func GenerateSessionID() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}

// HashPassword generates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// CompareHashAndPassword compares a bcrypt hashed password with its possible plaintext equivalent.
// Returns nil on success, or an error on failure.
func CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// IsAuthenticated checks if the user is logged in by checking if the session ID cookie exists in the database
// Returns true if the user is logged in, false otherwise. If an error occurs, it returns false.
// Used for rendering buttons based on the user's authentication status.
func IsAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("sessionId")
	if err != nil {
		return false
	}

	sessionQuery := "SELECT id FROM sessions WHERE id = ?"
	var id string

	stmt, err := db.GetDB().Prepare(sessionQuery)
	if err != nil {
		log.Println(err)
		return false
	}
	defer stmt.Close()

	err = stmt.QueryRow(cookie.Value).Scan(&id)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
