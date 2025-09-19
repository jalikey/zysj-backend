package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jalikey/zysj-backend/internal/repository" // !! 修改为你的模块路径
	"github.com/jalikey/zysj-backend/internal/models"
)
// Helper function to parse pagination (can be moved to a shared package)
func getPaginationParams(c *gin.Context) (page, limit, offset int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 { page = 1 }
	if limit < 1 { limit = 10 }
	if limit > 100 { limit = 100 }
	offset = (page - 1) * limit
	return
}
// GetCategories handles the GET request for retrieving all categories.
func GetCategories(c *gin.Context) {
	categories, err := repository.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
		return
	}

	// If no categories are found, return an empty array instead of null
	if categories == nil {
		categories = []models.Category{}
	}

	c.JSON(http.StatusOK, categories)
}
// ... GetCategories() 函数保持不变 ...

// GetArticlesByCategory handles getting all articles for a specific category.
func GetArticlesByCategory(c *gin.Context) {
	slug := c.Param("slug")

	category, err := repository.GetCategoryBySlug(slug)
	if err != nil {
		// ... error handling remains the same ...
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find category"})
		return
	}

	page, limit, offset := getPaginationParams(c)
	articles, totalItems, err := repository.GetArticlesByCategoryID(category.ID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve articles for this category"})
		return
	}
	
	// We wrap the original response in a new structure
	paginatedArticles := models.PaginatedResponse{
		Data: articles,
		Pagination: models.Pagination{
			CurrentPage: page,
			PageSize:    limit,
			TotalItems:  totalItems,
			TotalPages:  int(math.Ceil(float64(totalItems) / float64(limit))),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"category": category,
		"articles": paginatedArticles,
	})
}
