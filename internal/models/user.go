package models

import "time"

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Do not expose hash in JSON responses
	CreatedAt    time.Time `json:"created_at"`
}