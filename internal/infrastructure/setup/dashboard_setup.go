package setup

import (
	"belcamp/internal/infrastructure/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupDashboard(db *gorm.DB, protected *gin.RouterGroup) {
	h := &handlers.BaseHandler{}

	protected.GET("/", func(c *gin.Context) {
		dashboard(h, c)
	})
}

func dashboard(h *handlers.BaseHandler, c *gin.Context) {
	data := gin.H{
		"title":          "Dashboard",
		"totalOrders":    150,
		"totalProducts":  1250,
		"totalUsers":     250,
		"recentActivity": []string{},
	}

	h.Render(c, "dashboard.index", data, "")
}
