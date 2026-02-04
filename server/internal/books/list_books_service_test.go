//go:build medium

package books

import (
	"context"
	"testing"
	"time"

	"holocron/internal/books/domain"
)

func TestSearchBooksSource_WithKeyword_ReturnsMatchingBooks(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	_, _ = CreateBook(ctx, queries, CreateBookInput{
		Title:   "Go Programming",
		Authors: []string{"Author A"},
	})
	_, _ = CreateBook(ctx, queries, CreateBookInput{
		Title:   "Python Programming",
		Authors: []string{"Author B"},
	})

	query := "Go"
	keyword := domain.ToSearchKeyword(&query)
	pagination := domain.ToPagination(nil, nil)

	source := SearchBooksSource(queries)
	items, total, err := source(ctx, keyword, pagination)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if items[0].Title != "Go Programming" {
		t.Errorf("expected title 'Go Programming', got %s", items[0].Title)
	}
}

func TestSearchBooksSource_WithNilKeyword_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	pagination := domain.ToPagination(nil, nil)

	source := SearchBooksSource(queries)
	_, _, err := source(ctx, nil, pagination)

	if err != domain.ErrNotMyResponsibility {
		t.Errorf("expected ErrNotMyResponsibility, got %v", err)
	}
}

func TestListAllBooksSource_ReturnsAllBooks(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	_, _ = CreateBook(ctx, queries, CreateBookInput{
		Title:   "Book One",
		Authors: []string{"Author A"},
	})
	_, _ = CreateBook(ctx, queries, CreateBookInput{
		Title:   "Book Two",
		Authors: []string{"Author B"},
	})

	pagination := domain.ToPagination(nil, nil)

	source := ListAllBooksSource(queries)
	items, total, err := source(ctx, nil, pagination)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
}

func TestGetBookList_WithSources_UsesChainOfResponsibility(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	_, _ = CreateBook(ctx, queries, CreateBookInput{
		Title:   "Go Programming",
		Authors: []string{"Author A"},
	})
	_, _ = CreateBook(ctx, queries, CreateBookInput{
		Title:   "Python Programming",
		Authors: []string{"Author B"},
	})

	sources := []domain.BookListSource{
		SearchBooksSource(queries),
		ListAllBooksSource(queries),
	}

	// With keyword: SearchBooksSource handles
	query := "Go"
	keyword := domain.ToSearchKeyword(&query)
	pagination := domain.ToPagination(nil, nil)

	items, total, err := domain.GetBookList(ctx, sources, keyword, pagination)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || items[0].Title != "Go Programming" {
		t.Errorf("expected Go Programming, got %v", items)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}

	// Without keyword: ListAllBooksSource handles (fallback)
	items, total, err = domain.GetBookList(ctx, sources, nil, pagination)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
}

func TestGetBookList_WithDeletedBook_ExcludesDeletedBook(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	book1, _ := CreateBook(ctx, queries, CreateBookInput{
		Title:   "Book One",
		Authors: []string{"Author A"},
	})
	_, _ = CreateBook(ctx, queries, CreateBookInput{
		Title:   "Book Two",
		Authors: []string{"Author B"},
	})

	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, occurred_at)
		VALUES ('delete-event-1', ?, 'deleted', ?)
	`, book1.ID, book1.CreatedAt.Add(time.Second).UTC().Format(time.RFC3339))
	if err != nil {
		t.Fatalf("failed to insert delete event: %v", err)
	}

	sources := []domain.BookListSource{
		SearchBooksSource(queries),
		ListAllBooksSource(queries),
	}
	pagination := domain.ToPagination(nil, nil)

	items, total, err := domain.GetBookList(ctx, sources, nil, pagination)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if items[0].Title != "Book Two" {
		t.Errorf("expected title 'Book Two', got %s", items[0].Title)
	}
}
