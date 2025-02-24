package setup

import (
	"belcamp/internal/domain/entity"
	"belcamp/internal/infrastructure/persistence"
	"belcamp/internal/interfaces/http/handlers"
	"belcamp/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupCategories(db *gorm.DB, group *gin.RouterGroup) {
	// Create repository
	repo := persistence.NewGormRepository[entity.Category](db)

	// Add custom repository methods if needed
	categoryRepo := &persistence.CategoryRepository{
		GormRepository: repo.(*persistence.GormRepository[entity.Category]),
	}

	// Create service
	svc := service.NewCRUDService[entity.Category](categoryRepo)

	// Create handlers
	handler := handlers.NewCRUDHandler[entity.Category](svc, "categories")

	// Register routes
	handler.RegisterRoutes(group, "/categories")
}
