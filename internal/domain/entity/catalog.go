package entity

import "gorm.io/gorm"

// Catalog model
type Catalog struct {
	gorm.Model
	Name        string  `json:"name"`
	Slug        string  `gorm:"unique" json:"slug"`
	Description *string `json:"description,omitempty"`
	PDFPath     string  `json:"pdf_path"`
}
