// Package service provides the business logic for user operations
package service

import (
	"belcamp/internal/domain/entity"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService interface defines the methods for user operations
type AuthService interface {
	GetUserByID(id uint) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	ValidateCredentials(email, password string) (*entity.User, error)
}

// authService implements UserService
type authService struct {
	db *gorm.DB
}

// NewAuthService creates a new UserService instance
func NewAuthService(db *gorm.DB) AuthService {
	return &authService{
		db: db,
	}
}

// GetUserByID retrieves a user by their ID
func (s *authService) GetUserByID(id uint) (*entity.User, error) {
	var user entity.User
	result := s.db.Preload("Company.Address").First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email
func (s *authService) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	result := s.db.Where("email = ?", email).Preload("Company.Address").First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

// ValidateCredentials validates user credentials and returns the user if valid
func (s *authService) ValidateCredentials(email, password string) (*entity.User, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is approved
	if user.Status != "approved" {
		return nil, errors.New("account not approved")
	}

	return user, nil
}

// HashPassword creates a Laravel-compatible password hash
// func HashPassword(password string) (string, error) {
// 	// Laravel uses bcrypt with cost 12
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
// 	if err != nil {
// 		return "", err
// 	}

// 	return string(bytes), nil
// }

// CheckPassword verifies a password against a Laravel hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
