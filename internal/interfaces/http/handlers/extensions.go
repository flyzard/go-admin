// internal/interfaces/http/handlers/extensions.go
package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"belcamp/internal/domain/interfaces"
	"belcamp/internal/domain/valueobject"

	"github.com/gin-gonic/gin"
)

// List method extension for smart tables
func (h *CRUDHandler[T]) SmartTableList(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// Get pagination
	pagination := valueobject.NewPagination(page, pageSize)

	// Get entities
	entities, pagination, err := h.service.List(c.Request.Context(), pagination)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	// Get config - try to get from entity type first
	var config valueobject.SmartTableConfig

	// Create a zero value of T to check if it implements SmartTableProvider
	var zero T
	if provider, ok := interface{}(zero).(interfaces.SmartTableProvider); ok {
		config = provider.GetSmartTableConfig()
	} else {
		// Fallback to default config
		config = getDefaultConfig[T]()
	}

	// Apply sorting if provided
	sortField := c.DefaultQuery("sort", config.DefaultSort)
	sortOrder := c.DefaultQuery("order", config.DefaultOrder)

	if sortField != "" && len(entities) > 0 {
		sortEntities(entities, sortField, sortOrder == "desc")
	}

	// Build view model with config from entity
	viewModel := gin.H{
		"entities":        entities,
		"pagination":      pagination,
		"config":          config,
		"baseUrl":         c.Request.URL.Path,
		"currentSort":     sortField,
		"currentOrder":    sortOrder,
		"filter":          c.QueryMap("filter"),
		"currentPageSize": pageSize,
	}

	h.Render(c, h.tmpl+".index", viewModel, h.tmpl+".table")
}

// sortEntities sorts a slice of entities by a given field
func sortEntities[T any](entities []T, field string, descending bool) {
	sort.Slice(entities, func(i, j int) bool {
		// Get field values using reflection
		vI := reflect.ValueOf(entities[i])
		vJ := reflect.ValueOf(entities[j])

		// If field has dots, it's a nested field
		fields := strings.Split(field, ".")
		fI := getNestedField(vI, fields)
		fJ := getNestedField(vJ, fields)

		// Compare based on field type
		result := compare(fI, fJ)

		// Reverse if descending
		if descending {
			return !result
		}
		return result
	})
}

// getNestedField retrieves a nested field value using reflection
func getNestedField(v reflect.Value, fields []string) reflect.Value {
	current := v
	for _, field := range fields {
		// Handle pointers
		if current.Kind() == reflect.Ptr && !current.IsNil() {
			current = current.Elem()
		}

		// If we're working with a struct
		if current.Kind() == reflect.Struct {
			current = current.FieldByName(field)
		} else {
			// If not a struct, can't go further
			break
		}
	}
	return current
}

// compare compares two reflect values
func compare(a, b reflect.Value) bool {
	if !a.IsValid() || !b.IsValid() {
		return false
	}

	switch a.Kind() {
	case reflect.String:
		return a.String() < b.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a.Int() < b.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return a.Uint() < b.Uint()
	case reflect.Float32, reflect.Float64:
		return a.Float() < b.Float()
	case reflect.Bool:
		return !a.Bool() && b.Bool()
	default:
		return fmt.Sprint(a.Interface()) < fmt.Sprint(b.Interface())
	}
}

// getDefaultConfig creates a default configuration for entity type T
func getDefaultConfig[T any]() valueobject.SmartTableConfig {
	// Get type information
	t := reflect.TypeOf((*T)(nil)).Elem()

	// Create default columns based on struct fields
	var columns []valueobject.SmartTableColumn
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		// Skip complex fields like slices
		switch field.Type.Kind() {
		case reflect.Struct:
			if field.Anonymous {
				// Handle embedded structs like gorm.Model
				// For each field in the embedded struct, add a column
				embeddedType := field.Type
				for j := 0; j < embeddedType.NumField(); j++ {
					embeddedField := embeddedType.Field(j)
					if embeddedField.PkgPath == "" { // Exported
						columns = append(columns, valueobject.SmartTableColumn{
							Field:    embeddedField.Name,
							Label:    embeddedField.Name,
							Sortable: true,
							Visible:  true,
						})
					}
				}
			} else {
				columns = append(columns, valueobject.SmartTableColumn{
					Field:    field.Name,
					Label:    field.Name,
					Sortable: true,
					Visible:  true,
				})
			}
		case reflect.Slice, reflect.Map, reflect.Ptr:
			// Skip complex fields for default config
			continue
		default:
			columns = append(columns, valueobject.SmartTableColumn{
				Field:    field.Name,
				Label:    field.Name,
				Sortable: true,
				Visible:  true,
			})
		}
	}

	return valueobject.SmartTableConfig{
		Columns:      columns,
		DefaultSort:  "ID",
		DefaultOrder: "asc",
		PageSizes:    []int{10, 25, 50},
		Actions:      getDefaultActions(),
	}
}

// getDefaultActions creates standard CRUD actions
func getDefaultActions() []valueobject.SmartTableAction {
	return []valueobject.SmartTableAction{
		{
			Label:  "View",
			Action: "/{{.ID}}",
			Class:  "text-blue-600 hover:text-blue-900",
		},
		{
			Label:  "Edit",
			Action: "/{{.ID}}/edit",
			Class:  "text-green-600 hover:text-green-900",
		},
		{
			Label:   "Delete",
			Action:  "/{{.ID}}",
			Confirm: true,
			Message: "Are you sure you want to delete this item?",
			Class:   "text-red-600 hover:text-red-900",
		},
	}
}
