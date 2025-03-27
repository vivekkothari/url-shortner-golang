package repository

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"time"
	"url-shortner/db"
)

// URLRepository handles database interactions
type URLRepository struct {
	DB *db.Database
}

// NewURLRepository initializes the repository
func NewURLRepository(db *db.Database) *URLRepository {
	return &URLRepository{DB: db}
}

// SaveURL inserts a new URL or updates an existing one, returning the ID
func (r *URLRepository) SaveURL(ctx context.Context, longURL string) (int64, error) {
	query := `
		INSERT INTO url_shortner (id, long_url, created_at)
		VALUES ($1, $2, now())
		ON CONFLICT (long_url) DO UPDATE 
		SET long_url = excluded.long_url
		RETURNING id;
	`

	// Generate a random ID
	id := rand.Int63n(1<<63 - 1)

	var returnedID int64
	err := r.DB.Pool.QueryRow(ctx, query, id, longURL).Scan(&returnedID)
	if err != nil {
		return 0, errors.New("failed to insert/update URL: " + err.Error())
	}

	return returnedID, nil
}

// AccessLongURL updates last_accessed_at and retrieves long_url
func (r *URLRepository) AccessLongURL(ctx context.Context, id int64, now time.Time) (*string, error) {
	var longURL sql.NullString

	query := `
		UPDATE url_shortner 
		SET last_accessed_at = $1 
		WHERE id = $2 
		RETURNING long_url
	`

	err := r.DB.Pool.QueryRow(ctx, query, now, id).Scan(&longURL)
	if err != nil {
		return nil, err
	}

	if longURL.Valid {
		return &longURL.String, nil
	}
	return nil, nil // Return nil if no URL was found
}
