package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB a connection pool for the database
var DB *pgxpool.Pool

// ConnectDB establishes a connection to the PostgreSQL database
func ConnectDB() {
	var err error
	
	// Construct the database connection string from environment variables
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)

	// Create a new connection pool
	DB, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Ping the database to verify the connection
	if err := DB.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v\n", err)
	}

	fmt.Println("Successfully connected to the database!")
}

// CloseDB closes the database connection pool
func CloseDB() {
	if DB != nil {
		DB.Close()
		fmt.Println("Database connection closed.")
	}
}
