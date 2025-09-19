package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jalikey/zysj-backend/internal/models"
	"github.com/jalikey/zysj-backend/internal/repository"
)

type CategoryPayload struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
	ParentID    int64  `json:"parent_id"`
}

// CreateCategory handles POST requests to create a category.
func CreateCategory(c *gin.Context) {
	var payload CategoryPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{
		Name:        payload.Name,
		Slug:        payload.Slug,
		Description: payload.Description,
	}
	if payload.ParentID > 0 {
		category.ParentID = models.NullInt64{Int64: payload.ParentID, Valid: true}
	}

	_, err := repository.CreateCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully"})
}

// UpdateCategory handles PUT requests to update a category.
func UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}
	
	var payload CategoryPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{
		ID:          id,
		Name:        payload.Name,
		Slug:        payload.Slug,
		Description: payload.Description,
	}
	if payload.ParentID > 0 {
		category.ParentID = models.NullInt64{Int64: payload.ParentID, Valid: true}
	}

	if err := repository.UpdateCategory(category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}

// DeleteCategory handles DELETE requests to remove a category.
func DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	if err := repository.DeleteCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
// ... (之前的 CUD handlers) ...

// GetCategoryByID handles GET request for a single category by ID.
func GetCategoryByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	category, err := repository.GetCategoryByID(id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve category"})
		return
	}

	c.JSON(http.StatusOK, category)
}