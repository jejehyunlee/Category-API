package main

import (
	"Category-API/database"
	"Category-API/handlers"
	"Category-API/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug" // default
	}

	gin.SetMode(ginMode)

	if ginMode == "release" {
		log.Println("PRODUCTION mode activated")
	} else {
		log.Println("DEBUG mode activated")
	}

	// Initialize database
	database.ConnectDatabase()

	// Create router
	router := gin.New()

	// Custom logger yang skip /metrics
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Skip logging for /metrics endpoint
		if param.Path == "/metrics" {
			return ""
		}

		// Custom log format untuk endpoint lainnya
		return fmt.Sprintf("[GIN] %s | %3d | %13v | %15s | %-7s %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
		)
	}))

	router.SetTrustedProxies(nil)

	// Recovery middleware
	router.Use(gin.Recovery())

	// CORS middleware
	router.Use(middleware.CORS())

	// Health check routes
	router.GET("/health", handlers.HealthCheck)
	router.GET("/health/db", handlers.HealthCheckDB)

	// Metrics endpoint (simple version - no logs)
	router.GET("/metrics", func(c *gin.Context) {
		// Return minimal response tanpa log
		c.JSON(200, gin.H{
			"status":    "ok",
			"service":   "category-api",
			"timestamp": time.Now().Unix(),
		})
	})

	// Root route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Category API is running",
			"version": "1.0.0",
			"endpoints": map[string]string{
				"GET /":                  "API info",
				"GET /health":            "Basic health check",
				"GET /health/db":         "Database health check",
				"GET /metrics":           "Metrics endpoint",
				"GET /categories":        "Get all categories",
				"POST /categories":       "Create new category",
				"GET /categories/:id":    "Get category by ID",
				"PUT /categories/:id":    "Update category",
				"DELETE /categories/:id": "Delete category",
			},
		})
	})

	// Category routes
	categoryRoutes := router.Group("/categories")
	{
		categoryRoutes.GET("/", handlers.GetAllCategories)
		categoryRoutes.POST("/", handlers.CreateCategory)
		categoryRoutes.GET("/:id", handlers.GetCategoryByID)
		categoryRoutes.PUT("/:id", handlers.UpdateCategory)
		categoryRoutes.DELETE("/:id", handlers.DeleteCategory)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)

	router.Run(":" + port)
}
