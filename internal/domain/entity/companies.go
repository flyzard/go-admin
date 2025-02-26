package entity

import "gorm.io/gorm"

type Companies struct {
	gorm.Model
	Name *string `json:"name,omitempty"`
}
