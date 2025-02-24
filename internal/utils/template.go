// Package utils provides utility functions for the application.
package utils

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"belcamp/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

// TemplateData is the data that is passed to the HTML templates
type TemplateData struct {
	User        *models.User
	CurrentPage string
	CurrentYear int
	Csrf_token  string
	MenuItems   []MenuItem
}

type MenuItem struct {
	Name string
	URL  string
}

// NewTemplateData creates a new TemplateData struct
func NewTemplateData(c *gin.Context) gin.H {
	var menuItems = []MenuItem{
		{"Dashboard", "/"},
		{"Products", "/products"},
		{"Orders", "/orders"},
		{"Users", "/users"},
		{"Categories", "/categories"},
	}

	data := gin.H{}

	data["User"] = getCurrentUser(c)
	data["CurrentPage"] = getCurrentPage(c)
	data["CurrentYear"] = time.Now().Year()
	data["Csrf_token"] = csrf.Token(c.Request)
	data["MenuItems"] = menuItems

	return data
}

// getCurrentUser retrieves the current user from the session
func getCurrentUser(c *gin.Context) *models.User {
	user, exists := c.Get("user")
	if !exists {
		return nil
	}
	return user.(*models.User)
}

// getCurrentPage gets the current page from the request path
func getCurrentPage(c *gin.Context) string {
	path := c.Request.URL.Path
	// Remove leading slash and return first segment
	if len(path) > 1 {
		return path[1:]
	}
	return "dashboard"
}

func toInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", value)
	}
}

// SetupTemplateFunctions adds custom functions to the template engine
func setupTemplateFunctions(r *gin.Engine) {
	r.SetFuncMap(template.FuncMap{
		// Math functions
		"add": func(a, b int) int {
			return a + b
		},
		"subtract": func(a, b int) int {
			return a - b
		},
		"multiply": func(a, b int) int {
			return a * b
		},
		"divide": func(a, b int) float64 {
			return float64(a) / float64(b)
		},
		// Format functions
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"formatMoney": func(amount float64) string {
			return fmt.Sprintf("$%.2f", amount)
		},
		// Helper functions
		"isEven": func(n int) bool {
			return n%2 == 0
		},
		"inc": func(n int) int {
			return n + 1
		},
		"gt": func(a, b interface{}) bool {
			aInt, _ := toInt(a)
			bInt, _ := toInt(b)
			return aInt > bInt
		},
		"lt": func(a, b interface{}) bool {
			aInt, _ := toInt(a)
			bInt, _ := toInt(b)
			return aInt < bInt
		},
		"eq": func(a, b string) bool {
			return a == b
		},

		"formatCurrency": func(amount float64) string {
			return fmt.Sprintf("$%.2f", amount)
		},
		"index": func(obj interface{}, key string) interface{} {
			if obj == nil {
				return nil
			}

			val := reflect.ValueOf(obj)
			if !val.IsValid() {
				return nil
			}

			switch val.Kind() {
			case reflect.Map:
				mapVal := val.MapIndex(reflect.ValueOf(key))
				if !mapVal.IsValid() {
					return nil
				}
				return mapVal.Interface()
			case reflect.Struct:
				fieldVal := val.FieldByName(key)
				if !fieldVal.IsValid() {
					return nil
				}
				return fieldVal.Interface()
			default:
				return nil
			}
		},
		"getField": func(entity interface{}, field string) interface{} {
			// Use reflection to get field value
			val := reflect.ValueOf(entity)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			return val.FieldByName(field).Interface()
		},
		"equalAny": func(a, b interface{}) bool {
			// Convert to strings for comparison
			aStr := fmt.Sprintf("%v", a)
			bStr := fmt.Sprintf("%v", b)
			return aStr == bStr
		},
	})
}

func SetupTemplates(r *gin.Engine) {
	setupTemplateFunctions(r)
	// Create a new template and specify the functions
	tmpl := template.New("")
	tmpl.Funcs(r.FuncMap)

	// Find all template files
	files, err := filepath.Glob("templates/**/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// // Parse each template file and use its relative path as the template name
	registerTemplateFiles(files, tmpl)

	files, err = filepath.Glob("templates/pages/**/*.html")
	if err != nil {
		log.Fatal(err)
	}

	registerTemplateFiles(files, tmpl)

	// Set the template engine
	r.SetHTMLTemplate(tmpl)

	// Serve static files
	r.Static("/assets", "./assets")
	r.Static("/public", "./public")
}

func registerTemplateFiles(files []string, tmpl *template.Template) {
	for _, file := range files {
		// Read the file content
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read template %s: %v", file, err)
		}

		// Get the relative path from the templates directory
		name := strings.TrimPrefix(file, "templates/")
		name = strings.TrimSuffix(name, ".html")
		name = strings.ReplaceAll(name, "/", ".")
		log.Printf("Loading template: %s as %s", file, name)

		// Parse the template with its path as name
		_, err = tmpl.New(name).Parse(string(content))
		if err != nil {
			log.Fatalf("Failed to parse template %s: %v", name, err)
		}
	}
}

func GetCurrentPage(c *gin.Context) string {
	path := c.Request.URL.Path
	// Remove leading slash
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		return "dashboard"
	}

	// Get first segment of path
	if idx := strings.Index(path, "/"); idx != -1 {
		path = path[:idx]
	}

	return path
}
