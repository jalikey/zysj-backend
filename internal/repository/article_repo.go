package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/jalikey/zysj-backend/internal/database" // !! 修改为你的模块路径
	"github.com/jalikey/zysj-backend/internal/models"   // !! 修改为你的模块路径
)

// GetAllArticles queries the database and returns all articles.
// GetAllArticles now supports pagination.
// It returns a slice of articles for the current page and the total count of all articles.
func GetAllArticles(limit, offset int) ([]models.Article, int64, error) {
	// Query for the current page of articles
	query := `SELECT id, title, content, category_id, author, source, created_at, updated_at 
			  FROM articles 
			  ORDER BY created_at DESC
			  LIMIT $1 OFFSET $2`
	
	rows, err := database.DB.Query(context.Background(), query, limit, offset)
	if err != nil {
		log.Printf("Error querying paginated articles: %v\n", err)
		return nil, 0, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		// ... (scan logic remains the same as before) ...
		var article models.Article
		var categoryID sql.NullInt64
		if err := rows.Scan(&article.ID, &article.Title, &article.Content, &categoryID, &article.Author, &article.Source, &article.CreatedAt, &article.UpdatedAt); err != nil {
			log.Printf("Error scanning article row: %v\n", err)
			return nil, 0, err
		}
		if categoryID.Valid {
			article.CategoryID = models.NullInt64{Int64: categoryID.Int64, Valid: true}
		}
		articles = append(articles, article)
	}

	// Query for the total count of articles
	var totalItems int64
	countQuery := `SELECT COUNT(*) FROM articles`
	err = database.DB.QueryRow(context.Background(), countQuery).Scan(&totalItems)
	if err != nil {
		log.Printf("Error counting articles: %v\n", err)
		return nil, 0, err
	}

	return articles, totalItems, nil
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

// GetArticlesByCategoryID now supports pagination.
func GetArticlesByCategoryID(categoryID int64, limit, offset int) ([]models.Article, int64, error) {
	query := `SELECT id, title, content, category_id, author, source, created_at, updated_at 
			  FROM articles 
			  WHERE category_id = $1
			  ORDER BY created_at DESC
			  LIMIT $2 OFFSET $3`
			  
	rows, err := database.DB.Query(context.Background(), query, categoryID, limit, offset)
	// ... (scan logic is identical to GetAllArticles) ...
	if err != nil { log.Printf("Error querying articles by category ID: %v\n", err); return nil, 0, err }
	defer rows.Close()
	var articles []models.Article
	for rows.Next() {
		var article models.Article
		var catID sql.NullInt64
		if err := rows.Scan(&article.ID, &article.Title, &article.Content, &catID, &article.Author, &article.Source, &article.CreatedAt, &article.UpdatedAt); err != nil {
			log.Printf("Error scanning article row: %v\n", err); return nil, 0, err
		}
        if catID.Valid { article.CategoryID = models.NullInt64{Int64: catID.Int64, Valid: true} }
		articles = append(articles, article)
	}


	var totalItems int64
	countQuery := `SELECT COUNT(*) FROM articles WHERE category_id = $1`
	err = database.DB.QueryRow(context.Background(), countQuery, categoryID).Scan(&totalItems)
	if err != nil {
		log.Printf("Error counting articles for category: %v\n", err)
		return nil, 0, err
	}

	return articles, totalItems, nil
}
// ... 其他 import 和函数 ...

// SearchArticles performs a full-text search on the articles table.
// SearchArticles now supports pagination.
func SearchArticles(query string, limit, offset int) ([]models.Article, int64, error) {
	sqlQuery := `SELECT id, title, content, category_id, author, source, created_at, updated_at,
				 ts_rank(content_tsv, plainto_tsquery('simple', $1)) as rank
				 FROM articles
				 WHERE content_tsv @@ plainto_tsquery('simple', $1)
				 ORDER BY rank DESC
				 LIMIT $2 OFFSET $3`
	
	rows, err := database.DB.Query(context.Background(), sqlQuery, query, limit, offset)
	// ... (scan logic is identical to previous search) ...
	if err != nil { log.Printf("Error searching articles: %v\n", err); return nil, 0, err }
	defer rows.Close()
	var articles []models.Article
	for rows.Next() {
		var article models.Article
		var categoryID sql.NullInt64
		var rank float32
		if err := rows.Scan(&article.ID, &article.Title, &article.Content, &categoryID, &article.Author, &article.Source, &article.CreatedAt, &article.UpdatedAt, &rank); err != nil {
			log.Printf("Error scanning searched article row: %v\n", err); return nil, 0, err
		}
        if categoryID.Valid { article.CategoryID = models.NullInt64{Int64: categoryID.Int64, Valid: true} }
		articles = append(articles, article)
	}

	var totalItems int64
	countQuery := `SELECT COUNT(*) FROM articles WHERE content_tsv @@ plainto_tsquery('simple', $1)`
	err = database.DB.QueryRow(context.Background(), countQuery, query).Scan(&totalItems)
	if err != nil {
		log.Printf("Error counting search results: %v\n", err)
		return nil, 0, err
	}

	return articles, totalItems, nil
}
// ... (package and imports) ...

// --- CUD Functions for Admin ---

// CreateArticle inserts a new article into the database and returns its ID.
func CreateArticle(article models.Article) (int64, error) {
	query := `INSERT INTO articles (title, content, category_id, author, source) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var articleID int64
	
	// Use NullInt64 for nullable category_id
	var categoryID sql.NullInt64
	if article.CategoryID.Valid {
		categoryID.Int64 = article.CategoryID.Int64
		categoryID.Valid = true
	}

	err := database.DB.QueryRow(context.Background(), query,
		article.Title, article.Content, categoryID, article.Author, article.Source).Scan(&articleID)
	if err != nil {
		log.Printf("Error creating article: %v", err)
		return 0, err
	}
	return articleID, nil
}

// UpdateArticle updates an existing article in the database.
func UpdateArticle(article models.Article) error {
	query := `UPDATE articles 
			  SET title = $1, content = $2, category_id = $3, author = $4, source = $5, updated_at = now()
			  WHERE id = $6`
			  
	var categoryID sql.NullInt64
	if article.CategoryID.Valid {
		categoryID.Int64 = article.CategoryID.Int64
		categoryID.Valid = true
	}

	_, err := database.DB.Exec(context.Background(), query,
		article.Title, article.Content, categoryID, article.Author, article.Source, article.ID)
	if err != nil {
		log.Printf("Error updating article: %v", err)
	}
	return err
}

// DeleteArticle removes an article from the database by its ID.
func DeleteArticle(id int64) error {
	query := `DELETE FROM articles WHERE id = $1`
	_, err := database.DB.Exec(context.Background(), query, id)
	if err != nil {
		log.Printf("Error deleting article: %v", err)
	}
	return err
}

// ... (Existing Read functions like GetAllArticles, GetArticleByID etc. remain unchanged) ...