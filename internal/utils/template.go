// Package utils provides utility functions for the application.
package utils

import (
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
