//go:build small

package domain

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// When ParseBorrowingBook with any valid borrowedAt then returns correct BorrowedAt time
func TestParseBorrowingBook_WithAnyValidBorrowedAt_ReturnsCorrectTime(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns BorrowedAt matching the input time", prop.ForAll(
		func(borrowedAt time.Time) bool {
			borrowedAtStr := borrowedAt.UTC().Format(time.RFC3339)
			authorsJSON := "[]"
			book, err := ParseBorrowingBook("id", "title", &authorsJSON, borrowedAtStr, nil, nil, nil, nil, nil)
			if err != nil {
				return false
			}
			return book.BorrowedAt.Equal(borrowedAt.UTC().Truncate(time.Second))
		},
		gen.Time(),
	))
	properties.TestingRun(t)
}

// When ParseBorrowingBook with any valid dueDate then returns correct DueDate time
func TestParseBorrowingBook_WithAnyValidDueDate_ReturnsCorrectDueDate(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns DueDate matching the input time", prop.ForAll(
		func(dueDate time.Time) bool {
			dueDateStr := dueDate.UTC().Format(time.RFC3339)
			authorsJSON := "[]"
			book, err := ParseBorrowingBook("id", "title", &authorsJSON, "2024-01-15T10:30:00Z", nil, nil, nil, nil, &dueDateStr)
			if err != nil {
				return false
			}
			if book.DueDate == nil {
				return false
			}
			return book.DueDate.Equal(dueDate.UTC().Truncate(time.Second))
		},
		gen.Time(),
	))
	properties.TestingRun(t)
}

// When ParseBorrowingBook with any valid authors slice then returns correct Authors
func TestParseBorrowingBook_WithAnyValidAuthors_ReturnsCorrectAuthors(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns Authors matching the input slice", prop.ForAll(
		func(authors []string) bool {
			if authors == nil {
				authors = []string{}
			}
			authorsBytes, err := json.Marshal(authors)
			if err != nil {
				return false
			}
			authorsJSON := string(authorsBytes)
			book, err := ParseBorrowingBook("id", "title", &authorsJSON, "2024-01-15T10:30:00Z", nil, nil, nil, nil, nil)
			if err != nil {
				return false
			}
			return reflect.DeepEqual(book.Authors, authors)
		},
		gen.SliceOf(gen.AnyString()),
	))
	properties.TestingRun(t)
}

// When ParseBorrowingBook with nil optional fields then returns BorrowingBook with nil optional fields
func TestParseBorrowingBook_WithNilOptionalFields_ReturnsBorrowingBook(t *testing.T) {
	authorsJSON := `["著者A"]`

	book, err := ParseBorrowingBook("id", "title", &authorsJSON, "2024-01-15T10:30:00Z", nil, nil, nil, nil, nil)

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
	book, err := ParseBorrowingBook("id", "title", nil, "2024-01-15T10:30:00Z", nil, nil, nil, nil, nil)

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

	_, err := ParseBorrowingBook("id", "title", &invalidJSON, "2024-01-15T10:30:00Z", nil, nil, nil, nil, nil)

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
	invalidDueDate := "not-a-date"

	_, err := ParseBorrowingBook("id", "title", &authorsJSON, "2024-01-15T10:30:00Z", nil, nil, nil, nil, &invalidDueDate)

	if !errors.Is(err, ErrInvalidBorrowingBook) {
		t.Errorf("expected ErrInvalidBorrowingBook, got %v", err)
	}
}
