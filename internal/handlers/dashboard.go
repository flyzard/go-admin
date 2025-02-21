package handlers

import (
	"github.com/gin-gonic/gin"
)

// DashboardHandler handles dashboard requests
type DashboardHandler struct {
	BaseHandler
}

// NewDashboardHandler creates a new DashboardHandler
func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

// Dashboard renders the dashboard page
func (h *DashboardHandler) Dashboard(c *gin.Context) {
	data := gin.H{
		"title": "Dashboard",
		// Add any dashboard specific data here
	}

	h.Render(c, "pages/dashboard.html", data)
}

// Stats returns the dashboard stats
func (h *DashboardHandler) Stats(c *gin.Context) {
	// This is an example of an HTMX endpoint that returns a partial update
	data := gin.H{
		"totalOrders":  150,
		"totalUsers":   1250,
		"totalRevenue": 25000.50,
	}

	// If it's an HTMX request, return just the stats partial
	if c.GetHeader("HX-Request") == "true" {
		c.HTML(200, "partials/dashboard/stats.html", data)
		return
	}

	// Otherwise render the full page
	h.Render(c, "pages/dashboard.html", data)
}
