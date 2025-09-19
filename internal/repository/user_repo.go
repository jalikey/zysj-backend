package repository

import (
	"context"
	"log"

	"github.com/jalikey/zysj-backend/internal/database"
	"github.com/jalikey/zysj-backend/internal/models"
)

// CreateUser inserts a new user into the database.
func CreateUser(user models.User) (int64, error) {
	var userID int64
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id`
	err := database.DB.QueryRow(context.Background(), query, user.Username, user.PasswordHash).Scan(&userID)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return 0, err
	}
	return userID, nil
}

// GetUserByUsername finds a user by their username.
func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	query := `SELECT id, username, password_hash, created_at FROM users WHERE username = $1`
	err := database.DB.QueryRow(context.Background(), query, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt,
	)
	if err != nil {
		// It's common for this query to find no rows, which is not a server error
		return models.User{}, err
	}
	return user, nil
}