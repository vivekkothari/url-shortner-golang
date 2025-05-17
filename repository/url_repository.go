package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5"
	"math/rand"
	"time"
	"url-shortner/db"
	"url-shortner/internal/jetdb/url-shortner/public/model"
	"url-shortner/internal/jetdb/url-shortner/public/table"
)

var urlShortner = table.URLShortner

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
	// Generate a random ID
	id := rand.Int63n(1<<63 - 1)

	shortner := model.URLShortner{
		ID:        id,
		LongURL:   longURL,
		CreatedAt: time.Now(),
	}
	stmt := urlShortner.INSERT(urlShortner.ID, urlShortner.LongURL, urlShortner.CreatedAt).
		MODEL(shortner).ON_CONFLICT(urlShortner.LongURL).
		DO_UPDATE(postgres.SET(urlShortner.LongURL.SET(urlShortner.EXCLUDED.LongURL))).RETURNING(urlShortner.ID)

	var returnedID int64
	q, args := stmt.Sql()
	err := r.DB.Pool.QueryRow(ctx, q, args...).Scan(&returnedID)
	if err != nil {
		return 0, errors.New("failed to insert/update URL: " + err.Error())
	}

	return returnedID, nil
}

// AccessLongURL updates last_accessed_at and retrieves long_url
func (r *URLRepository) AccessLongURL(ctx context.Context, id int64, now time.Time) (string, error) {
	var longURL sql.NullString

	stmt := urlShortner.UPDATE(urlShortner.LastAccessedAt).SET(now).WHERE(urlShortner.ID.EQ(postgres.Int64(id))).
		RETURNING(urlShortner.LongURL)
	q, args := stmt.Sql()

	err := r.DB.Pool.QueryRow(ctx, q, args...).Scan(&longURL)
	if err != nil {
		return "", err
	}

	if longURL.Valid {
		return longURL.String, nil
	}
	return "", err // Return nil if no URL was found
}

func (r *URLRepository) getLongUrl(ctx context.Context, id int64) (*model.URLShortner, error) {
	statement := urlShortner.SELECT(urlShortner.AllColumns).
		FROM(urlShortner).
		WHERE(urlShortner.ID.EQ(postgres.Int64(id)))
	q, args := statement.Sql()
	row, _ := r.DB.Pool.Query(ctx, q, args...)
	return pgx.CollectOneRow(row, pgx.RowToAddrOfStructByPos[model.URLShortner])
}
