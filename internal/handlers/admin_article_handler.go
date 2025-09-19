package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jalikey/zysj-backend/internal/models"
	"github.com/jalikey/zysj-backend/internal/repository"
)

type ArticlePayload struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	CategoryID int64  `json:"category_id"`
	Author     string `json:"author"`
	Source     string `json:"source"`
}

// CreateArticle handles POST requests to create a new article.
func CreateArticle(c *gin.Context) {
	var payload ArticlePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article := models.Article{
		Title:   payload.Title,
		Content: payload.Content,
		Author:  payload.Author,
		Source:  payload.Source,
	}
	if payload.CategoryID > 0 {
		article.CategoryID = models.NullInt64{Int64: payload.CategoryID, Valid: true}
	}

	newID, err := repository.CreateArticle(article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
		return
	}

	createdArticle, _ := repository.GetArticleByID(newID)
	c.JSON(http.StatusCreated, createdArticle)
}

// UpdateArticle handles PUT requests to update an article.
func UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	var payload ArticlePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article := models.Article{
		ID:      id,
		Title:   payload.Title,
		Content: payload.Content,
		Author:  payload.Author,
		Source:  payload.Source,
	}
	if payload.CategoryID > 0 {
		article.CategoryID = models.NullInt64{Int64: payload.CategoryID, Valid: true}
	}

	if err := repository.UpdateArticle(article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article updated successfully"})
}

// DeleteArticle handles DELETE requests to remove an article.
func DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	if err := repository.DeleteArticle(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}