package handlers

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/jalikey/zysj-backend/internal/repository" // !! 修改为你的模块路径
	"github.com/jalikey/zysj-backend/internal/models"
)

// GetArticles handles the GET request for retrieving all articles.
func GetArticles(c *gin.Context) {
	articles, err := repository.GetAllArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve articles"})
		return
	}
	
	if articles == nil {
		articles = []repository.models.Article{}
	}

	c.JSON(http.StatusOK, articles)
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
	// Get search query from URL parameter ?q=...
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query cannot be empty"})
		return
	}

	articles, err := repository.SearchArticles(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform search"})
		return
	}

	if articles == nil {
		articles = []repository.models.Article{}
	}

	c.JSON(http.StatusOK, articles)
}