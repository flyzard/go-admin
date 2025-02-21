// Package models provides the models for the application.
package models

import (
	"database/sql"
	"time"
)

// Model base model for common fields
type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

// User represents the user model
type User struct {
	Model
	Name            string `gorm:"size:255;not null"`
	Email           string `gorm:"size:255;not null;unique"`
	Password        string `gorm:"size:255;not null"`
	Status          string `gorm:"size:50;default:'new'"`
	EmailVerifiedAt sql.NullTime
	CompanyID       *uint
	Company         *Company
	RememberToken   string `gorm:"size:100"`
}

// Company represents the company model
type Company struct {
	Model
	Name      string `gorm:"size:255;not null"`
	NIF       string `gorm:"size:255"`
	AddressID uint   `gorm:"not null"`
	Address   Address
	Phone     string `gorm:"size:255;not null"`
}

// Address represents the address model
type Address struct {
	Model
	Street  string `gorm:"size:255;not null"`
	ZipCode string `gorm:"size:255;not null"`
	City    string `gorm:"size:255;not null"`
	Town    string `gorm:"size:255;not null"`
	Country string `gorm:"size:255;not null"`
}

// Product represents the product model
type Product struct {
	Model
	Name             string `gorm:"size:255"`
	ShortDescription string `gorm:"size:255"`
	Description      string `gorm:"type:text"`
	Status           bool   `gorm:"default:true"`
	Slug             string `gorm:"size:255;not null"`
	CategoryID       *uint
	Category         *Category
	Datasheet        string `gorm:"size:255"`
	Prices           string `gorm:"type:json"`
	Measures         string `gorm:"type:json"`
	Photos           string `gorm:"type:json"`
	ColorPhotos      string `gorm:"type:json"`
	Sizes            string `gorm:"type:json"`
	Variants         []ProductVariant
}

// ProductVariant represents the product variant model
type ProductVariant struct {
	Model
	ProductID       uint   `gorm:"not null"`
	SKU             string `gorm:"size:20;not null"`
	Prices          string `gorm:"type:json"`
	Size            string `gorm:"size:20"`
	Availability    int    `gorm:"default:0"`
	Status          bool   `gorm:"default:true"`
	Colors          string `gorm:"type:json"`
	NextArrivalQty  *int
	NextArrivalDate *time.Time
}

// Category represents the category model
type Category struct {
	Model
	Name     string `gorm:"size:255;not null"`
	Slug     string `gorm:"size:255"`
	Icon     string `gorm:"size:255"`
	IsActive bool   `gorm:"default:true"`
	ParentID *uint
	Parent   *Category
	Order    int16 `gorm:"default:0"`
	InMenu   bool  `gorm:"default:true"`
	Products []Product
}
