package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `json:"name"`
	Email           string         `gorm:"uniqueIndex" json:"email"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at,omitempty"`
	Password        string         `json:"-"`                 // Hide from JSON
	RememberToken   *string        `gorm:"size:100" json:"-"` // Hide from JSON
	CompanyID       *uint          `json:"company_id,omitempty"`
	Status          string         `gorm:"type:enum('new','approved','rejected');default:new" json:"status"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`

	// Relations
	Company Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Orders  []Order `gorm:"foreignKey:UserID" json:"orders,omitempty"`
}
