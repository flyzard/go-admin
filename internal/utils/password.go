package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword creates a Laravel-compatible password hash
func HashPassword(password string) (string, error) {
	// Laravel uses bcrypt with cost 12
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// CheckPassword verifies a password against a Laravel hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
