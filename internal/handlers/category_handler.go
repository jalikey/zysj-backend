package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jalikey/zysj-backend/internal/repository" // !! 修改为你的模块路径
)

// GetCategories handles the GET request for retrieving all categories.
func GetCategories(c *gin.Context) {
	categories, err := repository.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
		return
	}

	// If no categories are found, return an empty array instead of null
	if categories == nil {
		categories = []repository.models.Category{}
	}

	c.JSON(http.StatusOK, categories)
}
// ... GetCategories() 函数保持不变 ...

// GetArticlesByCategory handles getting all articles for a specific category.
func GetArticlesByCategory(c *gin.Context) {
	slug := c.Param("slug")

	// First, find the category by slug to get its ID
	category, err := repository.GetCategoryBySlug(slug)
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find category"})
		return
	}

	// Then, fetch articles using the found category ID
	articles, err := repository.GetArticlesByCategoryID(category.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve articles for this category"})
		return
	}
    
    if articles == nil {
		articles = []repository.models.Article{}
	}

	// We can return the articles along with category info
	c.JSON(http.StatusOK, gin.H{
		"category": category,
		"articles": articles,
	})
}