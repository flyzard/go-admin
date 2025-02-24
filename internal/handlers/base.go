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

type MenuItem struct {
	Name string
	URL  string
}

var menuItems = []MenuItem{
	{"Dashboard", "/"},
	{"Products", "/products"},
	{"Orders", "/orders"},
	{"Users", "/users"},
}

// Render renders a template with the common template data
func (h *BaseHandler) Render(c *gin.Context, templateName string, data gin.H, partial string) {

	// Get base template data
	templateData := utils.NewTemplateData(c)

	// Merge template data with provided data
	if data == nil {
		data = gin.H{}
	}

	data["User"] = templateData.User
	data["CurrentYear"] = templateData.CurrentYear
	data["currentPage"] = utils.GetCurrentPage(c)
	data["csrf_token"] = csrf.Token(c.Request)
	data["MenuItems"] = menuItems

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
	c.HTML(status, "error.html", data)
}

func (h *BaseHandler) Redirect(c *gin.Context, path string) {
	// Handle the response based on request type
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", "/products")
		return
	}

	c.Redirect(http.StatusFound, "/products")
}
