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
			name TEXT,
			occurred_at TEXT NOT NULL
		);
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

// When GetBook with updated book then returns updated data
func TestGetBook_WithUpdatedBook_ReturnsUpdatedData(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	bookID := "test-book-id"
	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, code, title, authors, publisher, occurred_at)
		VALUES ('event-1', ?, 'created', 'CODE-001', 'Initial Title', '["Initial Author"]', 'Initial Publisher', '2024-01-01T00:00:00Z')
	`, bookID)
	if err != nil {
		t.Fatalf("failed to insert created event: %v", err)
	}

	expectedCode := "CODE-002"
	expectedTitle := "Updated Title"
	expectedAuthors := `["Updated Author 1", "Updated Author 2"]`
	expectedPublisher := "Updated Publisher"
	expectedPublishedDate := "2024-06-15"
	expectedThumbnailURL := "https://example.com/updated.jpg"

	_, err = db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, code, title, authors, publisher, published_date, thumbnail_url, occurred_at)
		VALUES ('event-2', ?, 'updated', ?, ?, ?, ?, ?, ?, '2024-01-02T00:00:00Z')
	`, bookID, expectedCode, expectedTitle, expectedAuthors, expectedPublisher, expectedPublishedDate, expectedThumbnailURL)
	if err != nil {
		t.Fatalf("failed to insert updated event: %v", err)
	}

	output, err := GetBook(ctx, queries, GetBookInput{BookID: bookID})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.ID != bookID {
		t.Errorf("expected ID %s, got %s", bookID, output.ID)
	}
	if output.Code == nil || *output.Code != expectedCode {
		t.Errorf("expected code %s, got %v", expectedCode, output.Code)
	}
	if output.Title != expectedTitle {
		t.Errorf("expected title %s, got %s", expectedTitle, output.Title)
	}
	if len(output.Authors) != 2 || output.Authors[0] != "Updated Author 1" || output.Authors[1] != "Updated Author 2" {
		t.Errorf("expected authors ['Updated Author 1', 'Updated Author 2'], got %v", output.Authors)
	}
	if output.Publisher == nil || *output.Publisher != expectedPublisher {
		t.Errorf("expected publisher %s, got %v", expectedPublisher, output.Publisher)
	}
	if output.PublishedDate == nil || *output.PublishedDate != expectedPublishedDate {
		t.Errorf("expected publishedDate %s, got %v", expectedPublishedDate, output.PublishedDate)
	}
	if output.ThumbnailURL == nil || *output.ThumbnailURL != expectedThumbnailURL {
		t.Errorf("expected thumbnailURL %s, got %v", expectedThumbnailURL, output.ThumbnailURL)
	}
}
