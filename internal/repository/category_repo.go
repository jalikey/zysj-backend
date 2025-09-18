package repository

import (
	"context"
	"log"
	"database/sql" // <--- 添加这一行
	"github.com/jalikey/zysj-backend/internal/database" // !! 修改为你的模块路径
	"github.com/jalikey/zysj-backend/internal/models"   // !! 修改为你的模块路径
)

// GetAllCategories queries the database and returns all categories.
func GetAllCategories() ([]models.Category, error) {
	query := `SELECT id, name, slug, description, parent_id, created_at FROM categories ORDER BY id ASC`

	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		log.Printf("Error querying categories: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		var category models.Category
		// For nullable parent_id, we need to scan into a sql.NullInt64 or similar
		var parentID sql.NullInt64 
		
		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&parentID, // Scan into the nullable type
			&category.CreatedAt,
		); err != nil {
			log.Printf("Error scanning category row: %v\n", err)
			return nil, err
		}
		
		// Convert sql.NullInt64 to our custom models.NullInt64
		if parentID.Valid {
			category.ParentID = models.NullInt64{Int64: parentID.Int64, Valid: true}
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error after iterating category rows: %v\n", err)
		return nil, err
	}

	return categories, nil
}
// ... GetAllCategories() 函数保持不变 ...

// GetCategoryBySlug queries for a single category by its slug.
func GetCategoryBySlug(slug string) (models.Category, error) {
	query := `SELECT id, name, slug, description, parent_id, created_at FROM categories WHERE slug = $1`
	var category models.Category
	var parentID sql.NullInt64

	row := database.DB.QueryRow(context.Background(), query, slug)
	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&parentID,
		&category.CreatedAt,
	)

	if err != nil {
		log.Printf("Error scanning single category row: %v\n", err)
		return models.Category{}, err
	}
    
    if parentID.Valid {
        category.ParentID = models.NullInt64{Int64: parentID.Int64, Valid: true}
    }

	return category, nil
}