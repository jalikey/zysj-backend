package models

import (
	"time"
)

// NullInt64 is an alias for sql.NullInt64 that can be used in the model.
// This handles the case where parent_id can be NULL in the database.
type NullInt64 struct {
	Int64 int64
	Valid bool // Valid is true if Int64 is not NULL
}

// Category represents the structure of our categories table
type Category struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	ParentID    NullInt64 `json:"parent_id,omitempty"` // Use NullInt64 for nullable foreign keys
	CreatedAt   time.Time `json:"created_at"`
}