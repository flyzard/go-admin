package entity

type ProductColorPhoto struct {
	ProductID uint       `gorm:"primaryKey" json:"product_id"`
	ColorID   uint       `gorm:"primaryKey" json:"color_id"`
	Photos    JSONPhotos `gorm:"type:json" json:"photos"`

	// Relations
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Color   Color   `gorm:"foreignKey:ColorID" json:"color,omitempty"`
}
