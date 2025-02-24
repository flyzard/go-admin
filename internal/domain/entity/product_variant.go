package entity

import (
	"time"

	"gorm.io/gorm"
)

type ProductVariant struct {
	gorm.Model
	ProductID       uint `gorm:"not null"`
	Product         *Product
	SKU             string `gorm:"size:20;not null"`
	Prices          string `gorm:"type:json"`
	Size            string `gorm:"size:20"`
	Availability    int    `gorm:"default:0"`
	Status          bool   `gorm:"default:true"`
	Colors          string `gorm:"type:json"`
	NextArrivalQty  *int
	NextArrivalDate *time.Time
}
