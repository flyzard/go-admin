package entity

import (
	"belcamp/internal/domain/valueobject"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type JSONMeasures map[string]interface{}
type JSONPhotos []string
type JSONColorPhotos map[string]interface{}
type JSONSizes []string
type JSONColors []map[string]interface{}

// Product model
type Product struct {
	ID               uint            `gorm:"primaryKey" json:"id"`
	Name             *string         `json:"name,omitempty"`
	ShortDescription *string         `json:"short_description,omitempty"`
	Description      *string         `json:"description,omitempty"`
	Status           bool            `gorm:"default:true" json:"status"`
	Slug             string          `json:"slug"`
	Prices           NullJSONPrices  `gorm:"type:json" json:"prices"`
	Measures         *JSONMeasures   `gorm:"type:json" json:"measures,omitempty"`
	Photos           *JSONPhotos     `gorm:"type:json" json:"photos,omitempty"`
	CategoryID       *uint           `json:"category_id,omitempty"`
	Datasheet        *string         `json:"datasheet,omitempty"`
	ColorPhotos      JSONColorPhotos `gorm:"type:json" json:"color_photos"`
	Sizes            *JSONSizes      `gorm:"type:json" json:"sizes,omitempty"`
	DeletedAt        gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`

	// Relations
	Category           *Category           `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	ProductVariants    []ProductVariant    `gorm:"foreignKey:ProductID" json:"product_variants,omitempty"`
	ProductColorPhotos []ProductColorPhoto `gorm:"foreignKey:ProductID" json:"product_color_photos,omitempty"`
}

// NullJSONPrices represents a nullable JSON field for prices
type NullJSONPrices struct {
	JSONPrices map[string]string
	Valid      bool // Valid is true if JSONPrices is not NULL
}

// Scan implements the Scanner interface for NullJSONPrices
func (nj *NullJSONPrices) Scan(value interface{}) error {
	if value == nil {
		nj.JSONPrices, nj.Valid = nil, false
		return nil
	}

	// Convert the value to []byte
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return errors.New(fmt.Sprintf("unsupported Scan, storing driver.Value type %T into type *NullJSONPrices", value))
	}

	// Unmarshal the JSON data
	var result map[string]string
	err := json.Unmarshal(data, &result)
	if err != nil {
		nj.Valid = false
		return err
	}

	// Set the result
	nj.JSONPrices = result
	nj.Valid = true
	return nil
}

// Value implements the driver.Valuer interface for NullJSONPrices
func (nj NullJSONPrices) Value() (driver.Value, error) {
	if !nj.Valid {
		return nil, nil
	}
	return json.Marshal(nj.JSONPrices)
}

// MarshalJSON implements the json.Marshaler interface for NullJSONPrices
func (nj NullJSONPrices) MarshalJSON() ([]byte, error) {
	if !nj.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nj.JSONPrices)
}

// UnmarshalJSON implements the json.Unmarshaler interface for NullJSONPrices
func (nj *NullJSONPrices) UnmarshalJSON(data []byte) error {
	// Check for null value
	if string(data) == "null" {
		nj.JSONPrices, nj.Valid = nil, false
		return nil
	}

	// Unmarshal into map
	var prices map[string]string
	if err := json.Unmarshal(data, &prices); err != nil {
		return err
	}

	nj.JSONPrices = prices
	nj.Valid = true
	return nil
}

func (j *JSONMeasures) Scan(value interface{}) error {
	if value == nil {
		*j = JSONMeasures{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to scan JSONMeasures: unexpected data type %T", value)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	*j = JSONMeasures(data)
	return nil
}

// Value implements the driver.Valuer interface for JSONMeasures
func (j JSONMeasures) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// MarshalJSON for JSONPhotos
func (j JSONPhotos) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(j))
}

// UnmarshalJSON for JSONPhotos
func (j *JSONPhotos) UnmarshalJSON(data []byte) error {
	var v []string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*j = JSONPhotos(v)
	return nil
}

// MarshalJSON for JSONSizes
func (j JSONSizes) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(j))
}

// UnmarshalJSON for JSONSizes
func (j *JSONSizes) UnmarshalJSON(data []byte) error {
	var v []string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*j = JSONSizes(v)
	return nil
}

// MarshalJSON for JSONColors
func (j JSONColors) MarshalJSON() ([]byte, error) {
	return json.Marshal([]map[string]interface{}(j))
}

// UnmarshalJSON for JSONColors
func (j *JSONColors) UnmarshalJSON(data []byte) error {
	var v []map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*j = JSONColors(v)
	return nil
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
				Field:      "GetCategoryName",
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
			// {
			// 	Label:   "Delete",
			// 	Icon:    "fas fa-trash",
			// 	Action:  "/products/{{.ID}}",
			// 	Confirm: true,
			// 	Message: "Are you sure you want to delete this product?",
			// 	Class:   "text-red-600 hover:text-red-900",
			// 	ShowWhen: func(entity interface{}) bool {
			// 		product, ok := entity.(Product)
			// 		if !ok {
			// 			return false
			// 		}
			// 		// Only show delete for products without variants
			// 		return len(product.Variants) == 0
			// 	},
			// },
		},
	}
}

func (p Product) GetCategoryName() string {
	if p.Category == nil {
		return ""
	}
	return p.Category.Name
}

func (p Product) FullPrice() (float64, error) {
	if !p.Prices.Valid {
		return 0, nil
	}

	prices, err := p.GetPricesMap()
	if err != nil {
		return 0, err
	}

	// Try to get the price for quantity 1
	if priceStr, exists := prices["1"]; exists {
		return strconv.ParseFloat(priceStr, 64)
	}

	return 0, nil
}

func (p Product) GetPricesMap() (map[string]string, error) {
	if !p.Prices.Valid {
		return map[string]string{}, nil
	}
	return p.Prices.JSONPrices, nil
}

// --- JSONPhotos ---

// Scan implements the sql.Scanner interface for JSONPhotos
func (j *JSONPhotos) Scan(value interface{}) error {
	if value == nil {
		*j = JSONPhotos{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to scan JSONPhotos: unexpected data type %T", value)
	}

	var data []string
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	*j = JSONPhotos(data)
	return nil
}

// Value implements the driver.Valuer interface for JSONPhotos
func (j JSONPhotos) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// --- JSONColorPhotos ---

// Scan implements the sql.Scanner interface for JSONColorPhotos
func (j *JSONColorPhotos) Scan(value interface{}) error {
	if value == nil {
		*j = JSONColorPhotos{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to scan JSONColorPhotos: unexpected data type %T", value)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	*j = JSONColorPhotos(data)
	return nil
}

// Value implements the driver.Valuer interface for JSONColorPhotos
func (j JSONColorPhotos) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// --- JSONSizes ---

// Scan implements the sql.Scanner interface for JSONSizes
func (j *JSONSizes) Scan(value interface{}) error {
	if value == nil {
		*j = JSONSizes{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to scan JSONSizes: unexpected data type %T", value)
	}

	var data []string
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	*j = JSONSizes(data)
	return nil
}

// Value implements the driver.Valuer interface for JSONSizes
func (j JSONSizes) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// --- JSONColors ---

// Scan implements the sql.Scanner interface for JSONColors
func (j *JSONColors) Scan(value interface{}) error {
	if value == nil {
		*j = JSONColors{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to scan JSONColors: unexpected data type %T", value)
	}

	var data []map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	*j = JSONColors(data)
	return nil
}

// Value implements the driver.Valuer interface for JSONColors
func (j JSONColors) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}
