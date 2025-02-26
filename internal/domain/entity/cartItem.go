package entity

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	CartID           uint    `json:"cart_id"`
	ProductVariantID uint    `json:"product_variant_id"`
	ProductID        *int    `json:"product_id,omitempty"`
	Quantity         int     `json:"quantity"`
	UnitPrice        float64 `gorm:"type:double(8,2);default:0.00" json:"unit_price"`

	// Relations
	Cart           Cart           `gorm:"foreignKey:CartID" json:"cart,omitempty"`
	ProductVariant ProductVariant `gorm:"foreignKey:ProductVariantID" json:"product_variant,omitempty"`
}
