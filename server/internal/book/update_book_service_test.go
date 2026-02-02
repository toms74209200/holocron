//go:build medium

package book

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"holocron/internal/book/domain"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// When UpdateBook with all fields then updates book
func TestUpdateBook_WithAllFields_UpdatesBook(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	bookID := uuid.New().String()
	initialTitle := uuid.New().String()
	initialAuthors := fmt.Sprintf(`["%s"]`, uuid.New().String())
	initialOccurredAt := time.Now().Add(-time.Duration(rand.Intn(720)+1) * time.Hour).Format(time.RFC3339)
	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, title, authors, occurred_at)
		VALUES (?, ?, 'created', ?, ?, ?)
	`, uuid.New().String(), bookID, initialTitle, initialAuthors, initialOccurredAt)
	if err != nil {
		t.Fatalf("failed to insert book: %v", err)
	}

	expectedCode := uuid.New().String()
	expectedTitle := uuid.New().String()
	expectedAuthors := []string{uuid.New().String()}
	expectedPublisher := uuid.New().String()
	expectedPublishedDate := fmt.Sprintf("%04d-%02d-%02d", 2020+rand.Intn(5), 1+rand.Intn(12), 1+rand.Intn(28))
	expectedThumbnailURL := "https://example.com/" + uuid.New().String() + ".jpg"

	output, err := UpdateBook(ctx, queries, UpdateBookInput{
		BookID:        bookID,
		Code:          &expectedCode,
		Title:         &expectedTitle,
		Authors:       &expectedAuthors,
		Publisher:     &expectedPublisher,
		PublishedDate: &expectedPublishedDate,
		ThumbnailURL:  &expectedThumbnailURL,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.ID != bookID {
		t.Errorf("expected ID %s, got %s", bookID, output.ID)
	}
	if output.Title != expectedTitle {
		t.Errorf("expected title %s, got %s", expectedTitle, output.Title)
	}
	if len(output.Authors) != len(expectedAuthors) || output.Authors[0] != expectedAuthors[0] {
		t.Errorf("expected authors %v, got %v", expectedAuthors, output.Authors)
	}
	if output.Code == nil || *output.Code != expectedCode {
		t.Errorf("expected code %s, got %v", expectedCode, output.Code)
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
	if output.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}

	var eventCount int
	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM book_events WHERE book_id = ? AND event_type = 'updated'`, bookID).Scan(&eventCount)
	if err != nil {
		t.Fatalf("failed to count events: %v", err)
	}
	expectedEventCount := 1
	if eventCount != expectedEventCount {
		t.Errorf("expected %d update event, got %d", expectedEventCount, eventCount)
	}
}

// When UpdateBook with partial fields then updates only specified fields
func TestUpdateBook_WithPartialFields_UpdatesOnlySpecifiedFields(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	bookID := uuid.New().String()
	initialCode := uuid.New().String()
	initialTitle := uuid.New().String()
	initialAuthor := uuid.New().String()
	initialAuthors := fmt.Sprintf(`["%s"]`, initialAuthor)
	initialPublisher := uuid.New().String()
	initialOccurredAt := time.Now().Add(-time.Duration(rand.Intn(720)+1) * time.Hour).Format(time.RFC3339)
	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, code, title, authors, publisher, occurred_at)
		VALUES (?, ?, 'created', ?, ?, ?, ?, ?)
	`, uuid.New().String(), bookID, initialCode, initialTitle, initialAuthors, initialPublisher, initialOccurredAt)
	if err != nil {
		t.Fatalf("failed to insert book: %v", err)
	}

	expectedTitle := uuid.New().String()
	output, err := UpdateBook(ctx, queries, UpdateBookInput{
		BookID: bookID,
		Title:  &expectedTitle,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Title != expectedTitle {
		t.Errorf("expected title %s, got %s", expectedTitle, output.Title)
	}
	if output.Code == nil || *output.Code != initialCode {
		t.Errorf("expected code to be preserved as %s, got %v", initialCode, output.Code)
	}
	expectedAuthors := []string{initialAuthor}
	if len(output.Authors) != len(expectedAuthors) || output.Authors[0] != expectedAuthors[0] {
		t.Errorf("expected authors to be preserved as %v, got %v", expectedAuthors, output.Authors)
	}
	if output.Publisher == nil || *output.Publisher != initialPublisher {
		t.Errorf("expected publisher to be preserved as %s, got %v", initialPublisher, output.Publisher)
	}
}

// When UpdateBook with non-existent book then returns not found error
func TestUpdateBook_WithNonExistentBook_ReturnsNotFoundError(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	nonExistentBookID := uuid.New().String()
	title := uuid.New().String()
	_, err := UpdateBook(ctx, queries, UpdateBookInput{
		BookID: nonExistentBookID,
		Title:  &title,
	})

	if !errors.Is(err, ErrBookNotFound) {
		t.Errorf("expected ErrBookNotFound, got %v", err)
	}
}

// When UpdateBook with no fields then creates update event with no changes
func TestUpdateBook_WithNoFields_CreatesNoChangeEvent(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	bookID := uuid.New().String()
	initialTitle := uuid.New().String()
	initialAuthor := uuid.New().String()
	initialAuthors := fmt.Sprintf(`["%s"]`, initialAuthor)
	initialOccurredAt := time.Now().Add(-time.Duration(rand.Intn(720)+1) * time.Hour).Format(time.RFC3339)
	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, title, authors, occurred_at)
		VALUES (?, ?, 'created', ?, ?, ?)
	`, uuid.New().String(), bookID, initialTitle, initialAuthors, initialOccurredAt)
	if err != nil {
		t.Fatalf("failed to insert book: %v", err)
	}

	output, err := UpdateBook(ctx, queries, UpdateBookInput{
		BookID: bookID,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Title != initialTitle {
		t.Errorf("expected title %s, got %s", initialTitle, output.Title)
	}

	var eventCount int
	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM book_events WHERE book_id = ? AND event_type = 'updated'`, bookID).Scan(&eventCount)
	if err != nil {
		t.Fatalf("failed to count events: %v", err)
	}
	expectedEventCount := 1
	if eventCount != expectedEventCount {
		t.Errorf("expected %d update event, got %d", expectedEventCount, eventCount)
	}
}

// When UpdateBook with invalid title then returns invalid title error
func TestUpdateBook_WithInvalidTitle_ReturnsInvalidTitleError(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	bookID := uuid.New().String()
	initialTitle := uuid.New().String()
	initialAuthors := fmt.Sprintf(`["%s"]`, uuid.New().String())
	initialOccurredAt := time.Now().Add(-time.Duration(rand.Intn(720)+1) * time.Hour).Format(time.RFC3339)
	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, title, authors, occurred_at)
		VALUES (?, ?, 'created', ?, ?, ?)
	`, uuid.New().String(), bookID, initialTitle, initialAuthors, initialOccurredAt)
	if err != nil {
		t.Fatalf("failed to insert book: %v", err)
	}

	tests := []struct {
		name  string
		title string
	}{
		{"empty title", ""},
		{"too long title", strings.Repeat("a", 201)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := UpdateBook(ctx, queries, UpdateBookInput{
				BookID: bookID,
				Title:  &tt.title,
			})

			if !errors.Is(err, domain.ErrInvalidTitle) {
				t.Errorf("expected ErrInvalidTitle, got %v", err)
			}
		})
	}
}

// When UpdateBook with invalid authors then returns invalid authors error
func TestUpdateBook_WithInvalidAuthors_ReturnsInvalidAuthorsError(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	bookID := uuid.New().String()
	initialTitle := uuid.New().String()
	initialAuthors := fmt.Sprintf(`["%s"]`, uuid.New().String())
	initialOccurredAt := time.Now().Add(-time.Duration(rand.Intn(720)+1) * time.Hour).Format(time.RFC3339)
	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, title, authors, occurred_at)
		VALUES (?, ?, 'created', ?, ?, ?)
	`, uuid.New().String(), bookID, initialTitle, initialAuthors, initialOccurredAt)
	if err != nil {
		t.Fatalf("failed to insert book: %v", err)
	}

	emptyAuthors := []string{}
	_, err = UpdateBook(ctx, queries, UpdateBookInput{
		BookID:  bookID,
		Authors: &emptyAuthors,
	})

	if !errors.Is(err, domain.ErrInvalidAuthors) {
		t.Errorf("expected ErrInvalidAuthors, got %v", err)
	}
}

// When UpdateBook multiple times then applies latest update
func TestUpdateBook_WithMultipleUpdates_AppliesLatestUpdate(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	ctx := context.Background()

	bookID := uuid.New().String()
	initialTitle := uuid.New().String()
	initialAuthors := fmt.Sprintf(`["%s"]`, uuid.New().String())
	initialOccurredAt := time.Now().Add(-time.Duration(rand.Intn(720)+1) * time.Hour).Format(time.RFC3339)
	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, title, authors, occurred_at)
		VALUES (?, ?, 'created', ?, ?, ?)
	`, uuid.New().String(), bookID, initialTitle, initialAuthors, initialOccurredAt)
	if err != nil {
		t.Fatalf("failed to insert book: %v", err)
	}

	firstTitle := uuid.New().String()
	_, err = UpdateBook(ctx, queries, UpdateBookInput{
		BookID: bookID,
		Title:  &firstTitle,
	})
	if err != nil {
		t.Fatalf("first update failed: %v", err)
	}

	secondTitle := uuid.New().String()
	output, err := UpdateBook(ctx, queries, UpdateBookInput{
		BookID: bookID,
		Title:  &secondTitle,
	})
	if err != nil {
		t.Fatalf("second update failed: %v", err)
	}

	if output.Title != secondTitle {
		t.Errorf("expected title %s, got %s", secondTitle, output.Title)
	}

	var eventCount int
	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM book_events WHERE book_id = ? AND event_type = 'updated'`, bookID).Scan(&eventCount)
	if err != nil {
		t.Fatalf("failed to count events: %v", err)
	}
	expectedEventCount := 2
	if eventCount != expectedEventCount {
		t.Errorf("expected %d update events, got %d", expectedEventCount, eventCount)
	}
}
