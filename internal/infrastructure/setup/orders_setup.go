package setup

import (
	"belcamp/internal/domain/entity"
	"belcamp/internal/infrastructure/handlers"
	"belcamp/internal/infrastructure/persistence"
	"belcamp/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupOrders(db *gorm.DB, group *gin.RouterGroup) {
	handlers.NewCRUDHandler(
		service.NewCRUDService(
			persistence.NewGormRepository[entity.Order](db),
		), "orders",
	).RegisterDefaultRoutes(group, "/orders")
}
