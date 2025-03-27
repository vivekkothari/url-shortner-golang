package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Database struct to hold the connection pool
type Database struct {
	Pool *pgxpool.Pool
}

// NewDatabase initializes the connection pool
func NewDatabase() (*Database, error) {
	dbURL := os.Getenv("DATABASE_URL") // Ex: "postgres://user:password@localhost:5432/mydb?sslmode=disable"

	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/url-shortner"
		//return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	// Connection pool config
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse DB config: %w", err)
	}
	config.MaxConns = 10                     // Set max connections
	config.MinConns = 2                      // Set min connections
	config.MaxConnIdleTime = 5 * time.Minute // Close idle connections

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create DB pool: %w", err)
	}

	// Verify the connection
	err = pool.Ping(context.Background())
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to connect to DB: %w", err)
	}

	log.Println("Connected to PostgreSQL successfully!")
	return &Database{Pool: pool}, nil
}

// Close closes the database connection pool
func (db *Database) Close() {
	db.Pool.Close()
	log.Println("Database connection pool closed")
}
