package setup

import (
	"belcamp/internal/infrastructure/handlers"
	"belcamp/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAuth(db *gorm.DB, public *gin.RouterGroup, protected *gin.RouterGroup) {
	h := handlers.NewAuthHandler(
		service.NewAuthService(db),
	)

	public.GET("/login", h.ShowLogin)
	public.POST("/login", h.Login)
	protected.POST("/logout", h.Logout)
}
