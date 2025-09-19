package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jalikey/zysj-backend/internal/models"
	"github.com/jalikey/zysj-backend/internal/repository"
)
// Helper function to parse pagination query parameters
func getPaginationParams(c *gin.Context) (page, limit, offset int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 { // Set a max limit
		limit = 100
	}

	offset = (page - 1) * limit
	return
}
// GetArticles handles the GET request for retrieving all articles.
func GetArticles(c *gin.Context) {
	page, limit, offset := getPaginationParams(c)

	articles, totalItems, err := repository.GetAllArticles(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve articles"})
		return
	}

	response := models.PaginatedResponse{
		Data: articles,
		Pagination: models.Pagination{
			CurrentPage: page,
			PageSize:    limit,
			TotalItems:  totalItems,
			TotalPages:  int(math.Ceil(float64(totalItems) / float64(limit))),
		},
	}

	c.JSON(http.StatusOK, response)
}
func GetArticleByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	article, err := repository.GetArticleByID(id)
	if err != nil {
		// pgx.ErrNoRows is the error for no result found
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve article"})
		return
	}

	c.JSON(http.StatusOK, article)
}
// ... 其他 import 和函数 ...

// SearchArticles handles the GET request for searching articles.
func SearchArticles(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query cannot be empty"})
		return
	}

	page, limit, offset := getPaginationParams(c)
	articles, totalItems, err := repository.SearchArticles(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform search"})
		return
	}

	response := models.PaginatedResponse{
		Data: articles,
		Pagination: models.Pagination{
			CurrentPage: page,
			PageSize:    limit,
			TotalItems:  totalItems,
			TotalPages:  int(math.Ceil(float64(totalItems) / float64(limit))),
		},
	}
	c.JSON(http.StatusOK, response)
}