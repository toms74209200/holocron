//go:build medium

package bookcode

import (
	"context"
	"database/sql"
	"testing"

	book "holocron/internal/book/domain"
	"holocron/internal/bookcode/domain"

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

func TestCreateBookByCode_WithValidCode_ReturnsBook(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	googleFetcher, err := NewGoogleBooksFetcher()
	if err != nil {
		t.Fatal(err)
	}
	openBDFetcher, err := NewOpenBDFetcher()
	if err != nil {
		t.Fatal(err)
	}
	sources := []domain.BookInfoSource{
		DBCacheSource(queries),
		ExternalAPISource(googleFetcher.Fetch, book.BookInfoFromGoogleBooks),
		ExternalAPISource(openBDFetcher.Fetch, book.BookInfoFromOpenBD),
	}

	output, err := CreateBookByCode(context.Background(), queries, sources, CreateBookByCodeInput{
		Code: "9784873115658",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Title != "リーダブルコード" {
		t.Errorf("expected title %q, got %q", "リーダブルコード", output.Title)
	}
}

func TestCreateBookByCode_WithEmptyCode_ReturnsInvalidCodeError(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)

	_, err := CreateBookByCode(context.Background(), queries, nil, CreateBookByCodeInput{
		Code: "",
	})

	if err != ErrInvalidCode {
		t.Errorf("expected ErrInvalidCode, got %v", err)
	}
}

func TestCreateBookByCode_WithSameCode_ReturnsBookWithSameInfo(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	googleFetcher, err := NewGoogleBooksFetcher()
	if err != nil {
		t.Fatal(err)
	}
	openBDFetcher, err := NewOpenBDFetcher()
	if err != nil {
		t.Fatal(err)
	}
	sources := []domain.BookInfoSource{
		DBCacheSource(queries),
		ExternalAPISource(googleFetcher.Fetch, book.BookInfoFromGoogleBooks),
		ExternalAPISource(openBDFetcher.Fetch, book.BookInfoFromOpenBD),
	}
	first, _ := CreateBookByCode(context.Background(), queries, sources, CreateBookByCodeInput{
		Code: "9784873115658",
	})

	second, err := CreateBookByCode(context.Background(), queries, sources, CreateBookByCodeInput{
		Code: "9784873115658",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if first.ID == second.ID {
		t.Error("expected different book IDs")
	}
	if second.Title != first.Title {
		t.Errorf("expected same title %q, got %q", first.Title, second.Title)
	}
}
