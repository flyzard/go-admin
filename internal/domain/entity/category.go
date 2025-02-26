package entity

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	// ID        uint           `gorm:"primaryKey" json:"id"`
	Name     string  `json:"name"`
	Slug     *string `json:"slug,omitempty"`
	Icon     *string `json:"icon,omitempty"`
	IsActive bool    `gorm:"default:true" json:"is_active"`
	ParentID *uint   `json:"parent_id,omitempty"`
	Order    *int16  `gorm:"default:0" json:"order,omitempty"`
	InMenu   bool    `gorm:"default:true" json:"in_menu"`
	// DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	// CreatedAt time.Time      `json:"created_at"`
	// UpdatedAt time.Time      `json:"updated_at"`

	// Relations
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Products []Product  `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}
