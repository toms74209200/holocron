//go:build medium

package book

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
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestGetBook_WithExistingBook_ReturnsBook(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	bookID := "test-book-id"
	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, title, authors, occurred_at)
		VALUES ('event-1', ?, 'created', 'Go入門', '["山田太郎"]', '2024-01-01T00:00:00Z')
	`, bookID)
	if err != nil {
		t.Fatalf("failed to insert book: %v", err)
	}

	output, err := GetBook(ctx, queries, GetBookInput{BookID: bookID})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.ID != bookID {
		t.Errorf("expected ID %s, got %s", bookID, output.ID)
	}
	if output.Title != "Go入門" {
		t.Errorf("expected title 'Go入門', got %s", output.Title)
	}
	if len(output.Authors) != 1 || output.Authors[0] != "山田太郎" {
		t.Errorf("expected authors ['山田太郎'], got %v", output.Authors)
	}
	if output.Status != "available" {
		t.Errorf("expected status 'available', got %s", output.Status)
	}
}

func TestGetBook_WithEmptyID_ReturnsInvalidIDError(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	_, err := GetBook(ctx, queries, GetBookInput{BookID: ""})

	if !errors.Is(err, ErrInvalidBookID) {
		t.Errorf("expected ErrInvalidBookID, got %v", err)
	}
}

func TestGetBook_WithNonExistentID_ReturnsNotFoundError(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	_, err := GetBook(ctx, queries, GetBookInput{BookID: "non-existent-id"})

	if !errors.Is(err, ErrBookNotFound) {
		t.Errorf("expected ErrBookNotFound, got %v", err)
	}
}

func TestGetBook_WithDeletedBook_ReturnsNotFoundError(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	bookID := "deleted-book-id"
	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, title, authors, occurred_at)
		VALUES ('event-1', ?, 'created', 'Deleted Book', '["Author"]', '2024-01-01T00:00:00Z')
	`, bookID)
	if err != nil {
		t.Fatalf("failed to insert book: %v", err)
	}
	_, err = db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, occurred_at)
		VALUES ('event-2', ?, 'deleted', '2024-01-02T00:00:00Z')
	`, bookID)
	if err != nil {
		t.Fatalf("failed to insert delete event: %v", err)
	}

	_, err = GetBook(ctx, queries, GetBookInput{BookID: bookID})

	if !errors.Is(err, ErrBookNotFound) {
		t.Errorf("expected ErrBookNotFound, got %v", err)
	}
}

func TestGetBook_WithNullOptionalFields_ReturnsNilPointers(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	bookID := "book-with-nulls"
	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, title, authors, occurred_at)
		VALUES ('event-1', ?, 'created', 'Minimal Book', '["Author"]', '2024-01-01T00:00:00Z')
	`, bookID)
	if err != nil {
		t.Fatalf("failed to insert book: %v", err)
	}

	output, err := GetBook(ctx, queries, GetBookInput{BookID: bookID})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Code != nil {
		t.Errorf("expected Code to be nil, got %v", output.Code)
	}
	if output.Publisher != nil {
		t.Errorf("expected Publisher to be nil, got %v", output.Publisher)
	}
	if output.ThumbnailURL != nil {
		t.Errorf("expected ThumbnailURL to be nil, got %v", output.ThumbnailURL)
	}
}

func TestGetBook_WithOptionalFields_ReturnsValues(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	bookID := "book-with-all-fields"
	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, code, title, authors, publisher, published_date, thumbnail_url, occurred_at)
		VALUES ('event-1', ?, 'created', '978-4-xxx', 'Full Book', '["Author"]', '技術評論社', '2024-01-01', 'https://example.com/thumb.jpg', '2024-01-01T00:00:00Z')
	`, bookID)
	if err != nil {
		t.Fatalf("failed to insert book: %v", err)
	}

	output, err := GetBook(ctx, queries, GetBookInput{BookID: bookID})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Code == nil || *output.Code != "978-4-xxx" {
		t.Errorf("expected Code '978-4-xxx', got %v", output.Code)
	}
	if output.Publisher == nil || *output.Publisher != "技術評論社" {
		t.Errorf("expected Publisher '技術評論社', got %v", output.Publisher)
	}
	if output.ThumbnailURL == nil || *output.ThumbnailURL != "https://example.com/thumb.jpg" {
		t.Errorf("expected ThumbnailURL, got %v", output.ThumbnailURL)
	}
}

func TestGetBook_WithReregisteredBook_ReturnsBook(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	bookID := "reregistered-book-id"
	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, title, authors, occurred_at)
		VALUES ('event-1', ?, 'created', 'Old Title', '["Old Author"]', '2024-01-01T00:00:00Z')
	`, bookID)
	if err != nil {
		t.Fatalf("failed to insert first created event: %v", err)
	}
	_, err = db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, occurred_at)
		VALUES ('event-2', ?, 'deleted', '2024-01-02T00:00:00Z')
	`, bookID)
	if err != nil {
		t.Fatalf("failed to insert deleted event: %v", err)
	}
	_, err = db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, title, authors, occurred_at)
		VALUES ('event-3', ?, 'created', 'New Title', '["New Author"]', '2024-01-03T00:00:00Z')
	`, bookID)
	if err != nil {
		t.Fatalf("failed to insert second created event: %v", err)
	}

	output, err := GetBook(ctx, queries, GetBookInput{BookID: bookID})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.ID != bookID {
		t.Errorf("expected ID %s, got %s", bookID, output.ID)
	}
	if output.Title != "New Title" {
		t.Errorf("expected title 'New Title', got %s", output.Title)
	}
	if len(output.Authors) != 1 || output.Authors[0] != "New Author" {
		t.Errorf("expected authors ['New Author'], got %v", output.Authors)
	}
}
