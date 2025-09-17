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
	
	apiV1 := router.Group("/api/v1")
	{
		// 新增搜索路由，放在前面
		apiV1.GET("/search", handlers.SearchArticles) 

		apiV1.GET("/categories", handlers.GetCategories)
		apiV1.GET("/categories/:slug", handlers.GetArticlesByCategory)
		
		apiV1.GET("/articles", handlers.GetArticles)
		apiV1.GET("/articles/:id", handlers.GetArticleByID)
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