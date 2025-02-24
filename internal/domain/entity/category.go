package entity

import "gorm.io/gorm"

// Category represents the category model
type Category struct {
	gorm.Model
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
