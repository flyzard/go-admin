package entity

type Color struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique" json:"name"`
	Code string `json:"code"`

	// Relations
	ProductColorPhotos []ProductColorPhoto `gorm:"foreignKey:ColorID" json:"product_color_photos,omitempty"`
}
