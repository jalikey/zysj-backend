package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jalikey/zysj-backend/internal/database" // !! 修改为你的模块路径
	"github.com/jalikey/zysj-backend/internal/handlers"   // !! 新增导入
)

func main() {
	// 1. Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables from OS")
	}

	// 2. Connect to the database
	database.ConnectDB()
	defer database.CloseDB()

	// 3. Initialize Gin router
	router := gin.Default()

	// 4. Setup routes
	// Simple health check route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	
	// Public API routes
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/login", handlers.Login)

		apiV1.GET("/search", handlers.SearchArticles)
		apiV1.GET("/categories", handlers.GetCategories)
		apiV1.GET("/categories/:slug", handlers.GetArticlesByCategory)
		// We keep the public GET routes for articles for simplicity
		apiV1.GET("/articles", handlers.GetArticles)
		apiV1.GET("/articles/:id", handlers.GetArticleByID)
	}

	// Admin API routes
	adminV1 := router.Group("/api/v1/admin")
	adminV1.Use(handlers.AuthMiddleware())
	{
		// Dashboard test route
		adminV1.GET("/dashboard", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Welcome to the admin dashboard!"})
		})

		// Articles CRUD
		// GET is already public, but we can add it here too if we want admin-specific logic later
		adminV1.GET("/articles", handlers.GetArticles) 
		adminV1.POST("/articles", handlers.CreateArticle)
		adminV1.PUT("/articles/:id", handlers.UpdateArticle)
		adminV1.DELETE("/articles/:id", handlers.DeleteArticle)

	// Categories CRUD
		adminV1.GET("/categories", handlers.GetCategories) 
		adminV1.GET("/categories/:id", handlers.GetCategoryByID) // 新增路由
		adminV1.POST("/categories", handlers.CreateCategory)
		adminV1.PUT("/categories/:id", handlers.UpdateCategory)
		adminV1.DELETE("/categories/:id", handlers.DeleteCategory)
	}

	// Admin API routes
	adminV1 := router.Group("/api/v1/admin")
	adminV1.Use(handlers.AuthMiddleware()) // Apply auth middleware to this group
	{
		// Test route to check if middleware works
		adminV1.GET("/dashboard", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Welcome to the admin dashboard!"})
		})

		// TODO: Add CRUD routes for articles and categories here in the next step
	}



	// 5. Start the server
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080" // Default port
	}

	log.Printf("Server starting on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}