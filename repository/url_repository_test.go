package repository

import (
	"context"
	"testing"
	"time"
	"url-shortner/db"

	_ "github.com/lib/pq"
)

func TestInsertAndGetURL(t *testing.T) {
	database, err := db.NewDatabase()
	repo := NewURLRepository(database)

	original := "https://example.com"

	id, err := repo.SaveURL(context.Background(), original)
	if err != nil {
		t.Fatalf("InsertURL failed: %v", err)
	}

	url, err := repo.getLongUrl(context.Background(), id)
	if err != nil {
		t.Fatalf("Unable to fetch url with id %d", id)
	}

	if url.LongURL != original {
		t.Fatalf("Unable to match %s and %s", url.LongURL, original)
	}

}

func TestInsertAndAccessUrl(t *testing.T) {
	database, err := db.NewDatabase()
	repo := NewURLRepository(database)

	original := "https://example123.com"

	id, err := repo.SaveURL(context.Background(), original)
	if err != nil {
		t.Fatalf("InsertURL failed: %v", err)
	}

	now := time.Now()

	longURL, err := repo.AccessLongURL(context.Background(), id, now)

	if err != nil {
		t.Fatalf("Unable to fetch url with id %d", id)
	}

	url, err := repo.getLongUrl(context.Background(), id)

	if err != nil {
		t.Fatalf("Unable to fetch url with id %d", id)
	}

	if !url.LastAccessedAt.Equal(now) {
		t.Fatalf("LastAccessedAt mismatch %s, %s ", url.LastAccessedAt, now)
	}

	if longURL != original {
		t.Fatalf("Unable to match %s and %s", longURL, original)
	}

}
