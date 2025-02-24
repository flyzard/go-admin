package entity

import (
	"belcamp/internal/domain/valueobject"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID               uint             `gorm:"primaryKey;column:id;autoIncrement"`
	Name             *string          `gorm:"column:name;type:varchar(255)"`
	ShortDescription *string          `gorm:"column:short_description;type:varchar(255)"`
	Description      *string          `gorm:"column:description;type:text"`
	Status           bool             `gorm:"column:status;type:tinyint(1);default:1;not null"`
	CreatedAt        *time.Time       `gorm:"column:created_at;type:timestamp"`
	UpdatedAt        *time.Time       `gorm:"column:updated_at;type:timestamp"`
	DeletedAt        gorm.DeletedAt   `gorm:"column:deleted_at;type:timestamp"`
	Slug             string           `gorm:"column:slug;type:varchar(255);not null"`
	Prices           json.RawMessage  `gorm:"column:prices;type:json;default:json_array();not null"`
	Measures         json.RawMessage  `gorm:"column:measures;type:json;default:json_array()"`
	Photos           json.RawMessage  `gorm:"column:photos;type:json"`
	CategoryID       *uint            `gorm:"column:category_id;foreignKey:categories(id)"`
	Datasheet        *string          `gorm:"column:datasheet;type:varchar(255)"`
	ColorPhotos      json.RawMessage  `gorm:"column:color_photos;type:json;default:json_array();not null"`
	Sizes            json.RawMessage  `gorm:"column:sizes;type:json;default:json_array()"`
	Variants         []ProductVariant `gorm:"foreignKey:product_id"`
}

// Price represents a price tier
type PriceTier struct {
	Quantity string `json:"quantity"`
	Price    string `json:"price"`
}

// GetPrices converts the JSON prices to a map
func (p *Product) GetPrices() (map[string]string, error) {
	prices := make(map[string]string)
	err := json.Unmarshal(p.Prices, &prices)
	return prices, err
}

// Measure represents product measurements
type Measurements struct {
	Weight         string `json:"Peso,omitempty"`
	Height         string `json:"Altura,omitempty"`
	Width          string `json:"Largura,omitempty"`
	Length         string `json:"Comprimento,omitempty"`
	BoxWeight      string `json:"Peso de caixa,omitempty"`
	PalletWeight   string `json:"Peso de palete,omitempty"`
	BoxHeight      string `json:"Altura de caixa,omitempty"`
	PalletHeight   string `json:"Altura de palete,omitempty"`
	BoxesPerPallet string `json:"Caixas em palete,omitempty"`
	BoxWidth       string `json:"Largura de caixa,omitempty"`
	QuantityPerBox string `json:"Quantidade em caixa,omitempty"`
	BoxLength      string `json:"Comprimento de caixa,omitempty"`
}

// GetMeasurements converts the JSON measures to a struct
func (p *Product) GetMeasurements() (*Measurements, error) {
	var measurements Measurements
	err := json.Unmarshal(p.Measures, &measurements)
	return &measurements, err
}

// GetPhotos returns the photos array
func (p *Product) GetPhotos() ([]string, error) {
	var photos []string
	err := json.Unmarshal(p.Photos, &photos)
	return photos, err
}

// ColorPhoto represents color and associated photos
type ColorInfo struct {
	Colors     string   `json:"colors"`
	Photos     []string `json:"photos"`
	ColorCodes []string `json:"color_codes"`
}

// GetColorPhotos returns the color photos map
func (p *Product) GetColorPhotos() (map[string]ColorInfo, error) {
	colorPhotos := make(map[string]ColorInfo)
	err := json.Unmarshal(p.ColorPhotos, &colorPhotos)
	return colorPhotos, err
}

// GetSizes returns the sizes array
func (p *Product) GetSizes() ([]string, error) {
	var sizes []string
	err := json.Unmarshal(p.Sizes, &sizes)
	return sizes, err
}

func (p Product) GetSmartTableConfig() valueobject.SmartTableConfig {
	return valueobject.SmartTableConfig{
		Columns: []valueobject.SmartTableColumn{
			{
				Field:      "Name",
				Label:      "Product Name",
				Sortable:   true,
				Filterable: true,
				FilterType: "text",
				Visible:    true,
			},
			{
				Field:      "ShortDescription",
				Label:      "Description",
				Sortable:   true,
				Filterable: true,
				FilterType: "text",
				Visible:    true,
			},
			{
				Field:      "Status",
				Label:      "Status",
				Sortable:   true,
				Filterable: true,
				FilterType: "select",
				FilterOpts: []valueobject.FilterOption{
					{Value: "true", Label: "Active"},
					{Value: "false", Label: "Inactive"},
				},
				Formatter: "formatStatus",
				Visible:   true,
			},
			{
				Field:      "Category.Name",
				Label:      "Category",
				Sortable:   true,
				Filterable: true,
				FilterType: "select",
				// This would be populated dynamically
				Visible: true,
			},
			{
				Field:     "CreatedAt",
				Label:     "Date Added",
				Sortable:  true,
				Formatter: "formatDate",
				Visible:   true,
			},
			{
				Field:     "Prices",
				Label:     "Price (Base)",
				Sortable:  false,
				Formatter: "formatPrice",
				Template:  "products/price_column.html", // Custom template for rendering prices
				Visible:   true,
			},
		},
		DefaultSort:  "Name",
		DefaultOrder: "asc",
		PageSizes:    []int{10, 25, 50, 100},
		Actions: []valueobject.SmartTableAction{
			{
				Label:  "View",
				Icon:   "fas fa-eye",
				Action: "/products/{{.ID}}",
				Class:  "text-blue-600 hover:text-blue-900",
			},
			{
				Label:  "Edit",
				Icon:   "fas fa-edit",
				Action: "/products/{{.ID}}/edit",
				Class:  "text-green-600 hover:text-green-900",
			},
			{
				Label:   "Delete",
				Icon:    "fas fa-trash",
				Action:  "/products/{{.ID}}",
				Confirm: true,
				Message: "Are you sure you want to delete this product?",
				Class:   "text-red-600 hover:text-red-900",
				ShowWhen: func(entity interface{}) bool {
					product, ok := entity.(Product)
					if !ok {
						return false
					}
					// Only show delete for products without variants
					return len(product.Variants) == 0
				},
			},
		},
	}
}
