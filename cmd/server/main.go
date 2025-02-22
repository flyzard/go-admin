// Package main provides the entry point for the application.
package main

import (
	"fmt"
	"log"
	"os"

	"belcamp/internal/database"
	"belcamp/internal/handlers"
	"belcamp/internal/middleware"
	"belcamp/internal/service"
	"belcamp/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
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

func initDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func setupRoutes(r *gin.Engine, db *gorm.DB) {
	userService := service.NewUserService(db)
	authHandler := handlers.NewAuthHandler(userService)

	// Public routes
	public := r.Group("/")
	public.Use(middleware.NoAuthMiddleware())
	{
		public.GET("/login", authHandler.ShowLogin)
		public.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		dashboardHandler := &handlers.DashboardHandler{}
		protected.GET("/", dashboardHandler.Dashboard)
		protected.POST("/logout", authHandler.Logout)

		// Product routes
		productService := service.NewProductService(db)
		productHandler := handlers.NewProductHandler(productService)
		protected.GET("/products", productHandler.List)
		protected.GET("/products/new", productHandler.ShowCreateForm)
		protected.POST("/products/new", productHandler.Create)
		protected.GET("/products/:id", productHandler.Show)
		protected.GET("/products/:id/edit", productHandler.ShowEditForm)
		protected.POST("/products/:id/edit", productHandler.Update)
		protected.POST("/products/:id/delete", productHandler.Delete)
	}

	// API routes group
	// api := r.Group("/api")
	// {
	// 	// API endpoints here
	// }
}
