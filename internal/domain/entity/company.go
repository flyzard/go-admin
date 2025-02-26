package entity

import (
	"time"

	"gorm.io/gorm"
)

type Company struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	NIF       *string        `json:"nif,omitempty"`
	Phone     string         `json:"phone"`
	AddressID uint           `json:"address_id"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	// Relations
	Address Address `gorm:"foreignKey:AddressID" json:"address,omitempty"`
	Carts   []Cart  `gorm:"foreignKey:CompanyID" json:"carts,omitempty"`
	Orders  []Order `gorm:"foreignKey:CompanyID" json:"orders,omitempty"`
	Users   []User  `gorm:"foreignKey:CompanyID" json:"users,omitempty"`
}
