// Package main provides the entry point for the application.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"belcamp/internal/database"
	"belcamp/internal/infrastructure/setup"
	"belcamp/internal/middleware"
	"belcamp/internal/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// Initialize configuration
	if err := initConfig(); err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get the underlying SQL DB connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}

	// Ensure the connection is closed when the application exits
	defer sqlDB.Close()

	// Initialize router
	router := initRouter()

	// Setup routes
	setupRoutes(router, db)

	// Start server with graceful shutdown
	startServer(router)
}

func initConfig() error {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	// You could expand this to initialize a proper config structure
	// config.Init() or similar if you need more sophisticated config handling
	return nil
}

func initRouter() *gin.Engine {
	// Set gin mode
	gin.SetMode(getEnv("GIN_MODE", "debug"))

	// Initialize Gin
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"}) // Trust the local proxy

	// Setup session middleware
	store := cookie.NewStore([]byte(getEnv("SESSION_SECRET", "your-secret-key")))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		Secure:   getEnv("GIN_MODE", "") == "release",
		HttpOnly: true,
	})
	r.Use(sessions.Sessions("belcamp_session", store))

	// Setup templates and middleware
	utils.SetupTemplates(r)
	r.Use(middleware.CSRF())

	return r
}

func setupRoutes(r *gin.Engine, db *gorm.DB) {
	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Dashboard routes
		setup.SetupDashboard(db, protected)

		// Product management
		setup.SetupProducts(db, protected)
		setup.SetupCategories(db, protected)

		// Order management
		setup.SetupOrders(db, protected)

		// User management
		setup.SetupUsers(db, protected)
	}

	// Public routes
	public := r.Group("/")
	public.Use(middleware.NoAuthMiddleware())
	{
		setup.SetupAuth(db, public, protected)
	}

	// API routes - uncomment and implement when needed
	// api := r.Group("/api/v1")
	// {
	//     // API endpoints here
	// }
}

func startServer(r *gin.Engine) {
	// Get port from environment
	port := getEnv("PORT", "8085")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

// Helper function to get environment variable with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
