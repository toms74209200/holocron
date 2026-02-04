//go:build medium

package books

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	_, err = db.Exec(`
		CREATE TABLE book_events (
			event_id TEXT PRIMARY KEY,
			book_id TEXT NOT NULL,
			event_type TEXT NOT NULL,
			code TEXT,
			title TEXT,
			authors TEXT,
			publisher TEXT,
			published_date TEXT,
			thumbnail_url TEXT,
			occurred_at TEXT NOT NULL
		);
		CREATE INDEX idx_book_events_book_id ON book_events(book_id);

		CREATE TABLE lending_events (
			event_id TEXT PRIMARY KEY,
			lending_id TEXT NOT NULL,
			book_id TEXT NOT NULL,
			borrower_id TEXT NOT NULL,
			event_type TEXT NOT NULL,
			due_date TEXT,
			occurred_at TEXT NOT NULL
		);
		CREATE INDEX idx_lending_events_lending_id ON lending_events(lending_id);
		CREATE INDEX idx_lending_events_book_id ON lending_events(book_id);
		CREATE INDEX idx_lending_events_borrower_id ON lending_events(borrower_id);

		CREATE TABLE user_events (
			event_id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			event_type TEXT NOT NULL,
			name TEXT NOT NULL,
			occurred_at TEXT NOT NULL
		);
		CREATE INDEX idx_user_events_user_id ON user_events(user_id);
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestCreateBook_WithValidInput_ReturnsOutput(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	input := CreateBookInput{
		Title:   "Test Book",
		Authors: []string{"Author1", "Author2"},
	}

	output, err := CreateBook(ctx, queries, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.ID == "" {
		t.Error("expected ID to be set")
	}
	if output.Title != "Test Book" {
		t.Errorf("expected Title 'Test Book', got %s", output.Title)
	}
	if len(output.Authors) != 2 {
		t.Errorf("expected 2 authors, got %d", len(output.Authors))
	}
	if output.Status != "available" {
		t.Errorf("expected Status 'available', got %s", output.Status)
	}
	if output.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestCreateBook_WithOptionalFields_ReturnsOutput(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	publisher := "Test Publisher"
	publishedDate := "2024-01-01"
	thumbnailURL := "https://example.com/thumb.jpg"

	input := CreateBookInput{
		Title:         "Test Book",
		Authors:       []string{"Author1"},
		Publisher:     &publisher,
		PublishedDate: &publishedDate,
		ThumbnailURL:  &thumbnailURL,
	}

	output, err := CreateBook(ctx, queries, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Publisher == nil || *output.Publisher != publisher {
		t.Errorf("expected Publisher %s, got %v", publisher, output.Publisher)
	}
	if output.PublishedDate == nil || *output.PublishedDate != publishedDate {
		t.Errorf("expected PublishedDate %s, got %v", publishedDate, output.PublishedDate)
	}
	if output.ThumbnailURL == nil || *output.ThumbnailURL != thumbnailURL {
		t.Errorf("expected ThumbnailURL %s, got %v", thumbnailURL, output.ThumbnailURL)
	}
}

func TestCreateBook_WithEmptyTitle_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	input := CreateBookInput{
		Title:   "",
		Authors: []string{"Author1"},
	}

	_, err := CreateBook(ctx, queries, input)

	if !errors.Is(err, ErrInvalidTitle) {
		t.Errorf("expected ErrInvalidTitle, got %v", err)
	}
}

func TestCreateBook_WithEmptyAuthors_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	input := CreateBookInput{
		Title:   "Test Book",
		Authors: []string{},
	}

	_, err := CreateBook(ctx, queries, input)

	if !errors.Is(err, ErrInvalidAuthors) {
		t.Errorf("expected ErrInvalidAuthors, got %v", err)
	}
}
