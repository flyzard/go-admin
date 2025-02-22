// Package handlers provides the handlers for the application.
package handlers

import (
	"net/http"

	"belcamp/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

// BaseHandler is a base handler for all handlers
type BaseHandler struct {
	// Add any common dependencies here
}

// Render renders a template with the common template data
func (h *BaseHandler) Render(c *gin.Context, templateName string, data gin.H) {

	// Get base template data
	templateData := utils.NewTemplateData(c)

	// Merge template data with provided data
	if data == nil {
		data = gin.H{}
	}

	data["User"] = templateData.User
	// data["CurrentPage"] = templateData.CurrentPage
	data["CurrentYear"] = templateData.CurrentYear
	data["currentPage"] = utils.GetCurrentPage(c)
	data["csrf_token"] = csrf.Token(c.Request)

	c.HTML(http.StatusOK, templateName, data)
}

// RenderError renders an error page
func (h *BaseHandler) RenderError(c *gin.Context, status int, message string) {
	data := gin.H{
		"error": message,
	}
	c.HTML(status, "error.html", data)
}
