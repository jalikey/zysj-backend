package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/jalikey/zysj-backend/internal/database" // !! 修改为你的模块路径
	"github.com/jalikey/zysj-backend/internal/models"   // !! 修改为你的模块路径
)

// GetAllArticles queries the database and returns all articles.
func GetAllArticles() ([]models.Article, error) {
	query := `
		SELECT id, title, content, category_id, author, source, created_at, updated_at 
		FROM articles 
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		log.Printf("Error querying articles: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article

	for rows.Next() {
		var article models.Article
		var categoryID sql.NullInt64

		if err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&categoryID,
			&article.Author,
			&article.Source,
			&article.CreatedAt,
			&article.UpdatedAt,
		); err != nil {
			log.Printf("Error scanning article row: %v\n", err)
			return nil, err
		}
		
		if categoryID.Valid {
			article.CategoryID = models.NullInt64{Int64: categoryID.Int64, Valid: true}
		}

		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error after iterating article rows: %v\n", err)
		return nil, err
	}

	return articles, nil
}


// ... GetAllArticles() 函数保持不变 ...

// GetArticleByID queries the database for a single article by its ID.
func GetArticleByID(id int64) (models.Article, error) {
	query := `
		SELECT id, title, content, category_id, author, source, created_at, updated_at 
		FROM articles 
		WHERE id = $1
	`
	var article models.Article
	var categoryID sql.NullInt64

	row := database.DB.QueryRow(context.Background(), query, id)
	err := row.Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&categoryID,
		&article.Author,
		&article.Source,
		&article.CreatedAt,
		&article.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error scanning single article row: %v\n", err)
		return models.Article{}, err
	}
	
	if categoryID.Valid {
		article.CategoryID = models.NullInt64{Int64: categoryID.Int64, Valid: true}
	}

	return article, nil
}

// GetArticlesByCategoryID queries articles belonging to a specific category.
func GetArticlesByCategoryID(categoryID int64) ([]models.Article, error) {
	query := `
		SELECT id, title, content, category_id, author, source, created_at, updated_at 
		FROM articles 
		WHERE category_id = $1
		ORDER BY created_at DESC
	`
	rows, err := database.DB.Query(context.Background(), query, categoryID)
	if err != nil {
		log.Printf("Error querying articles by category ID: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		var catID sql.NullInt64
		if err := rows.Scan(&article.ID, &article.Title, &article.Content, &catID, &article.Author, &article.Source, &article.CreatedAt, &article.UpdatedAt); err != nil {
			log.Printf("Error scanning article row: %v\n", err)
			return nil, err
		}
        if catID.Valid {
			article.CategoryID = models.NullInt64{Int64: catID.Int64, Valid: true}
		}
		articles = append(articles, article)
	}

	return articles, nil
}
// ... 其他 import 和函数 ...

// SearchArticles performs a full-text search on the articles table.
func SearchArticles(query string) ([]models.Article, error) {
	// Use plainto_tsquery for user-inputted search terms
	// It's safer and handles multiple words better
	sqlQuery := `
		SELECT id, title, content, category_id, author, source, created_at, updated_at,
		ts_rank(content_tsv, plainto_tsquery('simple', $1)) as rank
		FROM articles
		WHERE content_tsv @@ plainto_tsquery('simple', $1)
		ORDER BY rank DESC;
	`

	rows, err := database.DB.Query(context.Background(), sqlQuery, query)
	if err != nil {
		log.Printf("Error searching articles: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		var categoryID sql.NullInt64
		var rank float32 // To scan the rank value, though we might not use it in the struct

		if err := rows.Scan(
			&article.ID, &article.Title, &article.Content, &categoryID,
			&article.Author, &article.Source, &article.CreatedAt, &article.UpdatedAt,
			&rank,
		); err != nil {
			log.Printf("Error scanning searched article row: %v\n", err)
			return nil, err
		}
		
        if categoryID.Valid {
			article.CategoryID = models.NullInt64{Int64: categoryID.Int64, Valid: true}
		}
		articles = append(articles, article)
	}

	return articles, nil
}