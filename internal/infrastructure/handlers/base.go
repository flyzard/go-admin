// Package handlers provides the handlers for the application.
package handlers

import (
	"net/http"

	"belcamp/internal/utils"

	"github.com/gin-gonic/gin"
)

// BaseHandler is a base handler for all handlers
type BaseHandler struct {
	// Add any common dependencies here
}

type MenuItem struct {
	Name string
	URL  string
}

// Render renders a template with the common template data
func (h *BaseHandler) Render(c *gin.Context, templateName string, data gin.H, partial string) {

	// Get base template data
	templateData := utils.NewTemplateData(c)

	// Merge templateData into data
	for k, v := range templateData {
		// Only set if not already defined in data
		if _, exists := data[k]; !exists {
			data[k] = v
		}
	}

	if c.GetHeader("HX-Request") == "true" && partial != "" {
		c.HTML(http.StatusOK, partial, data)
		return
	}

	c.HTML(http.StatusOK, templateName, data)
}

// RenderError renders an error page
func (h *BaseHandler) RenderError(c *gin.Context, status int, message string) {
	data := gin.H{
		"error": message,
	}
	c.HTML(status, "error", data)
}

func (h *BaseHandler) Redirect(c *gin.Context, path string) {
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", path)
		c.Status(http.StatusFound)
		return
	}

	c.Redirect(http.StatusFound, path)
}
