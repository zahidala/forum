package utils

import (
	"forum/pkg/db"
	"log"
	"net/http"
	"time"

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

	sessionsQuery := "SELECT expiresAt FROM sessions WHERE id = ?"
	var expiresAtSession time.Time
	expiresAtCookie := cookie.Expires

	sessionStmt, sessionErr := db.GetDB().Prepare(sessionsQuery)
	if sessionErr != nil {
		log.Println(sessionErr)
		return false
	}
	defer sessionStmt.Close()

	sessionRowErr := sessionStmt.QueryRow(cookie.Value).Scan(&expiresAtSession)

	if sessionRowErr != nil {
		log.Println(err)
		return false
	}

	if time.Now().After(expiresAtSession) || time.Now().After(expiresAtCookie) {
		return false
	}

	return true
}
