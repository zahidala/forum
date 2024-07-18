package utils

import (
	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

func GenerateSessionID() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
