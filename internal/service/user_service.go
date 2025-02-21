// Package service provides the business logic for user operations
package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"belcamp/internal/models"
)

// UserService interface defines the methods for user operations
type UserService interface {
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
	ChangePassword(userID uint, currentPassword, newPassword string) error
	ValidateCredentials(email, password string) (*models.User, error)
	ListUsers(page, pageSize int) ([]models.User, int64, error)
	UpdateUserStatus(userID uint, status string) error
}

// userService implements UserService
type userService struct {
	db *gorm.DB
}

// NewUserService creates a new UserService instance
func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
}

// GetUserByID retrieves a user by their ID
func (s *userService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
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
func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := s.db.Where("email = ?", email).Preload("Company.Address").First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

// CreateUser creates a new user
func (s *userService) CreateUser(user *models.User) error {
	// Hash password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Set default status if not provided
	if user.Status == "" {
		user.Status = "new"
	}

	return s.db.Create(user).Error
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(user *models.User) error {
	// Don't update password through this method
	return s.db.Model(user).Omit("Password").Updates(user).Error
}

// DeleteUser soft deletes a user
func (s *userService) DeleteUser(id uint) error {
	return s.db.Delete(&models.User{}, id).Error
}

// ChangePassword changes a user's password
func (s *userService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash and save new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.db.Model(&user).Update("password", string(hashedPassword)).Error
}

// ValidateCredentials validates user credentials and returns the user if valid
func (s *userService) ValidateCredentials(email, password string) (*models.User, error) {
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

// ListUsers returns a paginated list of users
func (s *userService) ListUsers(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Get total count
	if err := s.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated users
	offset := (page - 1) * pageSize
	result := s.db.Preload("Company.Address").
		Offset(offset).
		Limit(pageSize).
		Find(&users)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, total, nil
}

// UpdateUserStatus updates a user's status
func (s *userService) UpdateUserStatus(userID uint, status string) error {
	// Validate status
	validStatuses := map[string]bool{
		"new":      true,
		"approved": true,
		"rejected": true,
	}

	if !validStatuses[status] {
		return errors.New("invalid status")
	}

	// Update user status
	result := s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
