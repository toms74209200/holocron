//go:build medium

package lending

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	_, err = db.Exec(`
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
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

type fakeBookQueries struct {
	countByBookId map[string]int64
}

func (f *fakeBookQueries) CountBookByBookId(_ context.Context, bookID string) (int64, error) {
	count, ok := f.countByBookId[bookID]
	if !ok {
		return 0, nil
	}
	return count, nil
}

// When BorrowBook with new book then returns output
func TestBorrowBook_WithNewBook_ReturnsOutput(t *testing.T) {
	db := setupTestDB(t)
	lendingQueries := New(db)
	bookQueries := &fakeBookQueries{
		countByBookId: map[string]int64{"book-123": 1},
	}
	ctx := context.Background()

	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	service := NewBorrowBookService(lendingQueries, bookQueries)
	service.now = func() time.Time { return fixedTime }

	input := BorrowBookInput{
		BookID:     "book-123",
		BorrowerID: "user-456",
		DueDays:    nil,
	}

	output, err := service.BorrowBook(ctx, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.ID == "" {
		t.Error("expected ID to be set")
	}
	if output.BookID != "book-123" {
		t.Errorf("expected BookID book-123, got %s", output.BookID)
	}
	if output.BorrowerID != "user-456" {
		t.Errorf("expected BorrowerID user-456, got %s", output.BorrowerID)
	}
	if !output.BorrowedAt.Equal(fixedTime) {
		t.Errorf("expected BorrowedAt %v, got %v", fixedTime, output.BorrowedAt)
	}
	expectedDueDate := fixedTime.AddDate(0, 0, 7)
	if !output.DueDate.Equal(expectedDueDate) {
		t.Errorf("expected DueDate %v, got %v", expectedDueDate, output.DueDate)
	}
}

// When BorrowBook with custom dueDays then returns output with custom due date
func TestBorrowBook_WithCustomDueDays_ReturnsOutputWithCustomDueDate(t *testing.T) {
	db := setupTestDB(t)
	lendingQueries := New(db)
	bookQueries := &fakeBookQueries{
		countByBookId: map[string]int64{"book-123": 1},
	}
	ctx := context.Background()

	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	service := NewBorrowBookService(lendingQueries, bookQueries)
	service.now = func() time.Time { return fixedTime }

	dueDays := 14
	input := BorrowBookInput{
		BookID:     "book-123",
		BorrowerID: "user-456",
		DueDays:    &dueDays,
	}

	output, err := service.BorrowBook(ctx, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedDueDate := fixedTime.AddDate(0, 0, 14)
	if !output.DueDate.Equal(expectedDueDate) {
		t.Errorf("expected DueDate %v, got %v", expectedDueDate, output.DueDate)
	}
}

// When BorrowBook with non-existent book then returns error
func TestBorrowBook_WithNonExistentBook_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	lendingQueries := New(db)
	bookQueries := &fakeBookQueries{
		countByBookId: map[string]int64{},
	}
	ctx := context.Background()

	service := NewBorrowBookService(lendingQueries, bookQueries)

	input := BorrowBookInput{
		BookID:     "non-existent",
		BorrowerID: "user-456",
		DueDays:    nil,
	}

	_, err := service.BorrowBook(ctx, input)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if err.Error() != "book not found" {
		t.Errorf("expected 'book not found' error, got %v", err)
	}
}

// When BorrowBook with already borrowed book by different user then returns error
func TestBorrowBook_WithAlreadyBorrowedBookByDifferentUser_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	lendingQueries := New(db)
	bookQueries := &fakeBookQueries{
		countByBookId: map[string]int64{"book-123": 1},
	}
	ctx := context.Background()

	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	service := NewBorrowBookService(lendingQueries, bookQueries)
	service.now = func() time.Time { return fixedTime }

	// First user borrows the book
	_, err := service.BorrowBook(ctx, BorrowBookInput{
		BookID:     "book-123",
		BorrowerID: "user-1",
		DueDays:    nil,
	})
	if err != nil {
		t.Fatalf("unexpected error on first borrow: %v", err)
	}

	// Second user tries to borrow the same book
	_, err = service.BorrowBook(ctx, BorrowBookInput{
		BookID:     "book-123",
		BorrowerID: "user-2",
		DueDays:    nil,
	})

	if !errors.Is(err, ErrBookAlreadyBorrowed) {
		t.Errorf("expected ErrBookAlreadyBorrowed, got %v", err)
	}
}

// When BorrowBook with already borrowed book by same user then extends due date
func TestBorrowBook_WithAlreadyBorrowedBookBySameUser_ExtendsDueDate(t *testing.T) {
	db := setupTestDB(t)
	lendingQueries := New(db)
	bookQueries := &fakeBookQueries{
		countByBookId: map[string]int64{"book-123": 1},
	}
	ctx := context.Background()

	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	service := NewBorrowBookService(lendingQueries, bookQueries)
	service.now = func() time.Time { return fixedTime }

	output1, err := service.BorrowBook(ctx, BorrowBookInput{
		BookID:     "book-123",
		BorrowerID: "user-1",
		DueDays:    nil,
	})
	if err != nil {
		t.Fatalf("unexpected error on first borrow: %v", err)
	}

	output2, err := service.BorrowBook(ctx, BorrowBookInput{
		BookID:     "book-123",
		BorrowerID: "user-1",
		DueDays:    nil,
	})
	if err != nil {
		t.Fatalf("unexpected error on extension: %v", err)
	}

	if output2.ID != output1.ID {
		t.Errorf("expected lending ID to remain %s, got %s", output1.ID, output2.ID)
	}

	expectedNewDueDate := output1.DueDate.AddDate(0, 0, 7)
	if !output2.DueDate.Equal(expectedNewDueDate) {
		t.Errorf("expected DueDate %v, got %v", expectedNewDueDate, output2.DueDate)
	}
}

// When BorrowBook with invalid dueDays then returns error
func TestBorrowBook_WithInvalidDueDays_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	lendingQueries := New(db)
	bookQueries := &fakeBookQueries{
		countByBookId: map[string]int64{"book-123": 1},
	}
	ctx := context.Background()

	service := NewBorrowBookService(lendingQueries, bookQueries)

	dueDays := 0
	input := BorrowBookInput{
		BookID:     "book-123",
		BorrowerID: "user-456",
		DueDays:    &dueDays,
	}

	_, err := service.BorrowBook(ctx, input)

	if err == nil {
		t.Error("expected error, got nil")
	}
}
