package setup

import (
	"belcamp/internal/domain/entity"
	"belcamp/internal/infrastructure/persistence"
	"belcamp/internal/interfaces/http/handlers"
	"belcamp/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupProducts(db *gorm.DB, group *gin.RouterGroup) {
	// Create repository
	repo := persistence.NewGormRepository[entity.Product](db)

	// Add custom repository methods if needed
	productRepo := &persistence.ProductRepository{
		GormRepository: repo.(*persistence.GormRepository[entity.Product]),
	}

	// Create service
	svc := service.NewCRUDService[entity.Product](productRepo)

	// Create handlers
	handler := handlers.NewCRUDHandler[entity.Product](svc, "products")

	// Register routes
	handler.RegisterRoutes(group, "/products")
}
