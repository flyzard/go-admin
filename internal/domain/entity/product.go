package entity

import (
	"belcamp/internal/domain/valueobject"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"gorm.io/gorm"
)

// Custom JSON types
type (
	JSONField       []byte
	JSONPrices      map[string]string
	JSONMeasures    map[string]any
	JSONPhotos      []string
	JSONColorPhotos map[string]any
	JSONSizes       []string
	JSONColors      []map[string]any
)

// Generic JSON field handling
func (j *JSONField) Scan(value any) error {
	if value == nil {
		*j = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*j = append((*j)[0:0], v...)
		return nil
	case string:
		*j = append((*j)[0:0], []byte(v)...)
		return nil
	default:
		return fmt.Errorf("unsupported scan, storing type %T into JSONField", value)
	}
}

func (j JSONField) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

// Product model
type Product struct {
	gorm.Model
	Name             *string   `json:"name,omitempty"`
	ShortDescription *string   `json:"short_description,omitempty"`
	Description      *string   `json:"description,omitempty"`
	Status           bool      `gorm:"default:true" json:"status"`
	Slug             string    `json:"slug"`
	Prices           JSONField `gorm:"type:json" json:"prices"`
	Measures         JSONField `gorm:"type:json" json:"measures,omitempty"`
	Photos           JSONField `gorm:"type:json" json:"photos,omitempty"`
	CategoryID       *uint     `json:"category_id,omitempty"`
	Datasheet        *string   `json:"datasheet,omitempty"`
	ColorPhotos      JSONField `gorm:"type:json" json:"color_photos"`
	Sizes            JSONField `gorm:"type:json" json:"sizes,omitempty"`

	// Relations
	Category           *Category           `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	ProductVariants    []ProductVariant    `gorm:"foreignKey:ProductID" json:"product_variants,omitempty"`
	ProductColorPhotos []ProductColorPhoto `gorm:"foreignKey:ProductID" json:"product_color_photos,omitempty"`

	// Cached values (not persisted)
	cachedPrices      *JSONPrices
	cachedMeasures    *JSONMeasures
	cachedPhotos      *JSONPhotos
	cachedColorPhotos *JSONColorPhotos
	cachedSizes       *JSONSizes
}

// JSON field getters
func (p *Product) GetPrices() (JSONPrices, error) {
	if p.cachedPrices != nil {
		return *p.cachedPrices, nil
	}

	var prices JSONPrices
	if len(p.Prices) == 0 {
		return prices, nil
	}

	if err := json.Unmarshal(p.Prices, &prices); err != nil {
		return nil, err
	}

	p.cachedPrices = &prices
	return prices, nil
}

func (p *Product) SetPrices(prices JSONPrices) error {
	data, err := json.Marshal(prices)
	if err != nil {
		return err
	}
	p.Prices = data
	p.cachedPrices = &prices
	return nil
}

func (p *Product) GetMeasures() (JSONMeasures, error) {
	if p.cachedMeasures != nil {
		return *p.cachedMeasures, nil
	}

	var measures JSONMeasures
	if len(p.Measures) == 0 {
		return measures, nil
	}

	if err := json.Unmarshal(p.Measures, &measures); err != nil {
		return nil, err
	}

	p.cachedMeasures = &measures
	return measures, nil
}

func (p *Product) SetMeasures(measures JSONMeasures) error {
	data, err := json.Marshal(measures)
	if err != nil {
		return err
	}
	p.Measures = data
	p.cachedMeasures = &measures
	return nil
}

func (p *Product) GetPhotos() (JSONPhotos, error) {
	if p.cachedPhotos != nil {
		return *p.cachedPhotos, nil
	}

	var photos JSONPhotos
	if len(p.Photos) == 0 {
		return photos, nil
	}

	if err := json.Unmarshal(p.Photos, &photos); err != nil {
		return nil, err
	}

	p.cachedPhotos = &photos
	return photos, nil
}

func (p *Product) SetPhotos(photos JSONPhotos) error {
	data, err := json.Marshal(photos)
	if err != nil {
		return err
	}
	p.Photos = data
	p.cachedPhotos = &photos
	return nil
}

func (p *Product) GetColorPhotos() (JSONColorPhotos, error) {
	if p.cachedColorPhotos != nil {
		return *p.cachedColorPhotos, nil
	}

	var colorPhotos JSONColorPhotos
	if len(p.ColorPhotos) == 0 {
		return colorPhotos, nil
	}

	if err := json.Unmarshal(p.ColorPhotos, &colorPhotos); err != nil {
		return nil, err
	}

	p.cachedColorPhotos = &colorPhotos
	return colorPhotos, nil
}

func (p *Product) SetColorPhotos(colorPhotos JSONColorPhotos) error {
	data, err := json.Marshal(colorPhotos)
	if err != nil {
		return err
	}
	p.ColorPhotos = data
	p.cachedColorPhotos = &colorPhotos
	return nil
}

func (p *Product) GetSizes() (JSONSizes, error) {
	if p.cachedSizes != nil {
		return *p.cachedSizes, nil
	}

	var sizes JSONSizes
	if len(p.Sizes) == 0 {
		return sizes, nil
	}

	if err := json.Unmarshal(p.Sizes, &sizes); err != nil {
		return nil, err
	}

	p.cachedSizes = &sizes
	return sizes, nil
}

func (p *Product) SetSizes(sizes JSONSizes) error {
	data, err := json.Marshal(sizes)
	if err != nil {
		return err
	}
	p.Sizes = data
	p.cachedSizes = &sizes
	return nil
}

// Business methods
func (p Product) CategoryName() string {
	if p.Category == nil {
		return ""
	}
	return p.Category.Name
}

func (p Product) FullPrice() (float64, error) {
	prices, err := p.GetPrices()
	if err != nil {
		return 0, err
	}

	// Try to get the price for quantity 1
	if priceStr, exists := prices["1"]; exists {
		return strconv.ParseFloat(priceStr, 64)
	}

	return 0, nil
}

func (p Product) MinimumPrice() (float64, error) {
	prices, err := p.GetPrices()
	if err != nil {
		return 0, err
	}

	minPrice := math.MaxFloat64
	for _, priceStr := range prices {
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			continue
		}
		if price < minPrice {
			minPrice = price
		}
	}

	if minPrice == math.MaxFloat64 {
		return 0, nil
	}

	return minPrice, nil
}

func (p Product) InStock() bool {
	if len(p.ProductVariants) == 0 {
		return false
	}

	for _, variant := range p.ProductVariants {
		if variant.Availability > 0 {
			return true
		}
	}

	return false
}

func (p Product) TotalStock() int {
	total := 0
	for _, variant := range p.ProductVariants {
		total += variant.Availability
	}
	return total
}

func (p Product) HasColor(colorName string) (bool, error) {
	colorPhotos, err := p.GetColorPhotos()
	if err != nil {
		return false, err
	}

	_, exists := colorPhotos[colorName]
	return exists, nil
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
				Visible: true,
			},
			{
				Field:      "CategoryName",
				Label:      "Category",
				Sortable:   true,
				Filterable: true,
				FilterType: "select",
				Visible:    true,
			},
			{
				Field:     "CreatedAt",
				Label:     "Date Added",
				Sortable:  true,
				Formatter: "formatDate",
				Visible:   true,
			},
			{
				Field:    "FullPrice",
				Label:    "Price (Base)",
				Sortable: false,
				Visible:  true,
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
		},
	}
}

// Custom JSON marshaling to handle computed fields
func (p Product) MarshalJSON() ([]byte, error) {
	type ProductAlias Product

	fullPrice, _ := p.FullPrice()
	minPrice, _ := p.MinimumPrice()

	return json.Marshal(&struct {
		ProductAlias
		FullPrice    float64 `json:"full_price"`
		MinimumPrice float64 `json:"minimum_price"`
		InStock      bool    `json:"in_stock"`
		TotalStock   int     `json:"total_stock"`
	}{
		ProductAlias: ProductAlias(p),
		FullPrice:    fullPrice,
		MinimumPrice: minPrice,
		InStock:      p.InStock(),
		TotalStock:   p.TotalStock(),
	})
}
