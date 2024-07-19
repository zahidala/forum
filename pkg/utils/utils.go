package utils

import (
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
