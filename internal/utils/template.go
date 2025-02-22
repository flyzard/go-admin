// Package utils provides utility functions for the application.
package utils

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"belcamp/internal/models"

	"github.com/gin-gonic/gin"
)

// TemplateData is the data that is passed to the HTML templates
type TemplateData struct {
	User        *models.User
	CurrentPage string
	CurrentYear int
}

// NewTemplateData creates a new TemplateData struct
func NewTemplateData(c *gin.Context) TemplateData {
	return TemplateData{
		User:        getCurrentUser(c),
		CurrentPage: getCurrentPage(c),
		CurrentYear: time.Now().Year(),
	}
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
func SetupTemplateFunctions(r *gin.Engine) {
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
	})
}

func SetupTemplates(r *gin.Engine) {
	// Create a new template and specify the functions
	tmpl := template.New("")
	tmpl.Funcs(r.FuncMap)

	// Find all template files
	files, err := filepath.Glob("internal/templates/**/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// // Parse each template file and use its relative path as the template name
	registerTemplateFiles(files, tmpl)

	files, err = filepath.Glob("internal/templates/pages/**/*.html")
	if err != nil {
		log.Fatal(err)
	}

	registerTemplateFiles(files, tmpl)

	// r.LoadHTMLGlob("internal/templates/**/*")

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
		name := strings.TrimPrefix(file, "internal/templates/")
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
