package entity

import "gorm.io/gorm"

type Cart struct {
	// ID        uint           `gorm:"primaryKey" json:"id"`
	gorm.Model
	CompanyID uint  `json:"company_id"`
	Status    int16 `gorm:"default:0" json:"status"`
	// DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	// CreatedAt time.Time      `json:"created_at"`
	// UpdatedAt time.Time      `json:"updated_at"`

	// Relations
	Company   Company    `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	CartItems []CartItem `gorm:"foreignKey:CartID" json:"cart_items,omitempty"`
	Orders    []Order    `gorm:"foreignKey:CartID" json:"orders,omitempty"`
}
