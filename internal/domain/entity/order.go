package entity

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	Status       int16    `gorm:"default:0" json:"status"`
	CartID       uint     `json:"cart_id"`
	UserID       uint     `json:"user_id"`
	CompanyID    uint     `json:"company_id"`
	IP           string   `json:"ip"`
	Notes        *string  `json:"notes,omitempty"`
	ShippingCost *float64 `gorm:"type:double(8,2)" json:"shipping_cost,omitempty"`
	Total        float64  `gorm:"type:double(8,2)" json:"total"`
	Withdraw     bool     `gorm:"default:false" json:"withdraw"`
	Taxes        float64  `gorm:"type:double(8,2);default:0.00" json:"taxes"`
	Weight       *float64 `gorm:"type:double(8,2)" json:"weight,omitempty"`

	// Relations
	Cart    Cart    `gorm:"foreignKey:CartID" json:"cart,omitempty"`
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Company Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}
