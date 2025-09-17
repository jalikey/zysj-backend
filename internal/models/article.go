package models

import "time"

// Article represents the structure of our articles table
type Article struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CategoryID NullInt64 `json:"category_id,omitempty"` // A category might be optional
	Author     string    `json:"author,omitempty"`
	Source     string    `json:"source,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}