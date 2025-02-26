// Package main provides the entry point for the application.
package main

import (
	"log"
	"os"

	"belcamp/internal/database"
	"belcamp/internal/infrastructure/setup"
	"belcamp/internal/middleware"

	"belcamp/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database connection
	db, err := database.Initialize()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Gin
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"}) // Trust the local proxy
	gin.SetMode(os.Getenv("GIN_MODE"))

	// Setup session middleware
	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		Secure:   os.Getenv("GIN_MODE") == "release",
		HttpOnly: true,
	})
	r.Use(sessions.Sessions("belcamp_session", store))

	utils.SetupTemplates(r)
	r.Use(middleware.CSRF())

	// Setup routes
	setupRoutes(r, db)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}

func setupRoutes(r *gin.Engine, db *gorm.DB) {
	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		setup.SetupDashboard(db, protected)
		setup.SetupProducts(db, protected)
		setup.SetupCategories(db, protected)
		setup.SetupOrders(db, protected)
		setup.SetupUsers(db, protected)
	}

	// Public routes
	public := r.Group("/")
	public.Use(middleware.NoAuthMiddleware())
	{
		setup.SetupAuth(db, public, protected)
	}

	// API routes group
	// api := r.Group("/api")
	// {
	// 	// API endpoints here
	// }
}
