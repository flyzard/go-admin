package entity

import (
	"time"

	"gorm.io/gorm"
)

type ProductVariant struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	ProductID       uint           `json:"product_id"`
	SKU             string         `gorm:"size:20" json:"sku"`
	Prices          JSONPrices     `gorm:"type:json" json:"prices"`
	Size            *string        `gorm:"size:20" json:"size,omitempty"`
	Availability    int            `gorm:"default:0" json:"availability"`
	Status          bool           `gorm:"default:true" json:"status"`
	Colors          *JSONColors    `gorm:"type:json" json:"colors,omitempty"`
	NextArrivalQty  *int           `json:"next_arrival_qty,omitempty"`
	NextArrivalDate *time.Time     `json:"next_arrival_date,omitempty"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`

	// Relations
	Product   Product    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	CartItems []CartItem `gorm:"foreignKey:ProductVariantID" json:"cart_items,omitempty"`
}
