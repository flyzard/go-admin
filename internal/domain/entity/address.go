package entity

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Street    string         `json:"street"`
	Zipcode   string         `json:"zipcode"`
	City      string         `json:"city"`
	Town      string         `json:"town"`
	Country   string         `json:"country"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	// Relations
	Companies []Company `gorm:"foreignKey:AddressID" json:"companies,omitempty"`
}
