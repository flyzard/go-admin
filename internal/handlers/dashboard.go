package handlers

import (
	"github.com/gin-gonic/gin"
)

// DashboardHandler handles dashboard requests
type DashboardHandler struct {
	BaseHandler
}

// Dashboard renders the dashboard page
func (h *DashboardHandler) Dashboard(c *gin.Context) {
	data := gin.H{
		"title":          "Dashboard",
		"totalOrders":    150,
		"totalProducts":  1250,
		"totalUsers":     250,
		"recentActivity": []string{},
	}

	h.Render(c, "dashboard.index", data)
}
