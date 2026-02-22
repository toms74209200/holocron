//go:build small

package domain

import (
	"errors"
	"testing"
	"time"
)

// When ParseBorrowingBook with valid input then returns BorrowingBook with correct values
func TestParseBorrowingBook_WithValidInput_ReturnsBorrowingBook(t *testing.T) {
	id := "550e8400-e29b-41d4-a716-446655440001"
	title := "テスト書籍"
	authorsJSON := `["著者A","著者B"]`
	borrowedAt := "2024-01-15T10:30:00Z"
	dueDate := "2024-01-22T10:30:00Z"
	code := "9784873119045"
	publisher := "オライリージャパン"

	book, err := ParseBorrowingBook(id, title, &authorsJSON, borrowedAt, &code, &publisher, nil, nil, &dueDate)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if book.ID != id {
		t.Errorf("expected ID %s, got %s", id, book.ID)
	}
	if book.Title != title {
		t.Errorf("expected Title %s, got %s", title, book.Title)
	}
	if len(book.Authors) != 2 || book.Authors[0] != "著者A" || book.Authors[1] != "著者B" {
		t.Errorf("expected Authors [著者A 著者B], got %v", book.Authors)
	}
	expectedBorrowedAt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	if !book.BorrowedAt.Equal(expectedBorrowedAt) {
		t.Errorf("expected BorrowedAt %v, got %v", expectedBorrowedAt, book.BorrowedAt)
	}
	if book.DueDate == nil {
		t.Fatal("expected DueDate to be set")
	}
	expectedDueDate := time.Date(2024, 1, 22, 10, 30, 0, 0, time.UTC)
	if !book.DueDate.Equal(expectedDueDate) {
		t.Errorf("expected DueDate %v, got %v", expectedDueDate, *book.DueDate)
	}
	if book.Code == nil || *book.Code != code {
		t.Errorf("expected Code %s, got %v", code, book.Code)
	}
	if book.Publisher == nil || *book.Publisher != publisher {
		t.Errorf("expected Publisher %s, got %v", publisher, book.Publisher)
	}
}

// When ParseBorrowingBook with nil optional fields then returns BorrowingBook with nil optional fields
func TestParseBorrowingBook_WithNilOptionalFields_ReturnsBorrowingBook(t *testing.T) {
	id := "550e8400-e29b-41d4-a716-446655440001"
	authorsJSON := `["著者A"]`
	borrowedAt := "2024-01-15T10:30:00Z"

	book, err := ParseBorrowingBook(id, "", &authorsJSON, borrowedAt, nil, nil, nil, nil, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if book.Code != nil {
		t.Errorf("expected Code to be nil, got %v", book.Code)
	}
	if book.Publisher != nil {
		t.Errorf("expected Publisher to be nil, got %v", book.Publisher)
	}
	if book.PublishedDate != nil {
		t.Errorf("expected PublishedDate to be nil, got %v", book.PublishedDate)
	}
	if book.ThumbnailURL != nil {
		t.Errorf("expected ThumbnailURL to be nil, got %v", book.ThumbnailURL)
	}
	if book.DueDate != nil {
		t.Errorf("expected DueDate to be nil, got %v", book.DueDate)
	}
}

// When ParseBorrowingBook with nil authorsJSON then returns BorrowingBook with empty authors
func TestParseBorrowingBook_WithNilAuthorsJSON_ReturnsEmptyAuthors(t *testing.T) {
	borrowedAt := "2024-01-15T10:30:00Z"

	book, err := ParseBorrowingBook("id", "title", nil, borrowedAt, nil, nil, nil, nil, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(book.Authors) != 0 {
		t.Errorf("expected empty authors, got %v", book.Authors)
	}
}

// When ParseBorrowingBook with invalid authorsJSON then returns error
func TestParseBorrowingBook_WithInvalidAuthorsJSON_ReturnsError(t *testing.T) {
	invalidJSON := "not-json"
	borrowedAt := "2024-01-15T10:30:00Z"

	_, err := ParseBorrowingBook("id", "title", &invalidJSON, borrowedAt, nil, nil, nil, nil, nil)

	if !errors.Is(err, ErrInvalidBorrowingBook) {
		t.Errorf("expected ErrInvalidBorrowingBook, got %v", err)
	}
}

// When ParseBorrowingBook with invalid borrowedAt then returns error
func TestParseBorrowingBook_WithInvalidBorrowedAt_ReturnsError(t *testing.T) {
	authorsJSON := `[]`

	_, err := ParseBorrowingBook("id", "title", &authorsJSON, "not-a-date", nil, nil, nil, nil, nil)

	if !errors.Is(err, ErrInvalidBorrowingBook) {
		t.Errorf("expected ErrInvalidBorrowingBook, got %v", err)
	}
}

// When ParseBorrowingBook with invalid dueDate then returns error
func TestParseBorrowingBook_WithInvalidDueDate_ReturnsError(t *testing.T) {
	authorsJSON := `[]`
	borrowedAt := "2024-01-15T10:30:00Z"
	invalidDueDate := "not-a-date"

	_, err := ParseBorrowingBook("id", "title", &authorsJSON, borrowedAt, nil, nil, nil, nil, &invalidDueDate)

	if !errors.Is(err, ErrInvalidBorrowingBook) {
		t.Errorf("expected ErrInvalidBorrowingBook, got %v", err)
	}
}
