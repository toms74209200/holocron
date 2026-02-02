//go:build medium

package lending

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
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
	bookID := uuid.New().String()
	userID := uuid.New().String()
	bookQueries := &fakeBookQueries{
		countByBookId: map[string]int64{bookID: 1},
	}
	ctx := context.Background()

	_, err := lendingQueries.GetCurrentLending(ctx, bookID)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("precondition failed: expected no current lending, got error: %v", err)
	}

	borrowTime := time.Now().Add(-time.Duration(rand.Intn(720)+1) * time.Hour)
	service := NewBorrowBookService(lendingQueries, bookQueries)
	service.now = func() time.Time { return borrowTime }

	input := BorrowBookInput{
		BookID:     bookID,
		BorrowerID: userID,
		DueDays:    nil,
	}

	output, err := service.BorrowBook(ctx, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.ID == "" {
		t.Error("expected ID to be set")
	}
	if output.BookID != bookID {
		t.Errorf("expected BookID %s, got %s", bookID, output.BookID)
	}
	if output.BorrowerID != userID {
		t.Errorf("expected BorrowerID %s, got %s", userID, output.BorrowerID)
	}
	if !output.BorrowedAt.Equal(borrowTime) {
		t.Errorf("expected BorrowedAt %v, got %v", borrowTime, output.BorrowedAt)
	}
	expectedDueDate := borrowTime.AddDate(0, 0, 7)
	if !output.DueDate.Equal(expectedDueDate) {
		t.Errorf("expected DueDate %v, got %v", expectedDueDate, output.DueDate)
	}

	currentLending, err := lendingQueries.GetCurrentLending(ctx, bookID)
	if err != nil {
		t.Fatalf("postcondition failed: expected current lending, got error: %v", err)
	}
	if currentLending.BookID != bookID {
		t.Errorf("postcondition failed: expected book ID %s, got %s", bookID, currentLending.BookID)
	}
	if currentLending.BorrowerID != userID {
		t.Errorf("postcondition failed: expected borrower ID %s, got %s", userID, currentLending.BorrowerID)
	}
}

// When BorrowBook with custom dueDays then returns output with custom due date
func TestBorrowBook_WithCustomDueDays_ReturnsOutputWithCustomDueDate(t *testing.T) {
	db := setupTestDB(t)
	lendingQueries := New(db)
	bookID := uuid.New().String()
	userID := uuid.New().String()
	bookQueries := &fakeBookQueries{
		countByBookId: map[string]int64{bookID: 1},
	}
	ctx := context.Background()

	borrowTime := time.Now().Add(-time.Duration(rand.Intn(720)+1) * time.Hour)
	service := NewBorrowBookService(lendingQueries, bookQueries)
	service.now = func() time.Time { return borrowTime }

	dueDays := rand.Intn(30) + 1
	input := BorrowBookInput{
		BookID:     bookID,
		BorrowerID: userID,
		DueDays:    &dueDays,
	}

	output, err := service.BorrowBook(ctx, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedDueDate := borrowTime.AddDate(0, 0, dueDays)
	if !output.DueDate.Equal(expectedDueDate) {
		t.Errorf("expected DueDate %v, got %v", expectedDueDate, output.DueDate)
	}
}

// When BorrowBook with non-existent book then returns error
func TestBorrowBook_WithNonExistentBook_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	lendingQueries := New(db)
	nonExistentBookID := uuid.New().String()
	userID := uuid.New().String()
	bookQueries := &fakeBookQueries{
		countByBookId: map[string]int64{},
	}
	ctx := context.Background()

	count, err := bookQueries.CountBookByBookId(ctx, nonExistentBookID)
	if err != nil {
		t.Fatalf("failed to verify precondition: %v", err)
	}
	if count != 0 {
		t.Fatalf("precondition failed: expected book count 0, got %d", count)
	}

	service := NewBorrowBookService(lendingQueries, bookQueries)

	input := BorrowBookInput{
		BookID:     nonExistentBookID,
		BorrowerID: userID,
		DueDays:    nil,
	}

	_, err = service.BorrowBook(ctx, input)

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
	bookID := uuid.New().String()
	firstUser := uuid.New().String()
	secondUser := uuid.New().String()
	bookQueries := &fakeBookQueries{
		countByBookId: map[string]int64{bookID: 1},
	}
	ctx := context.Background()

	borrowTime := time.Now().Add(-time.Duration(rand.Intn(720)+1) * time.Hour)
	service := NewBorrowBookService(lendingQueries, bookQueries)
	service.now = func() time.Time { return borrowTime }

	firstBorrow, err := service.BorrowBook(ctx, BorrowBookInput{
		BookID:     bookID,
		BorrowerID: firstUser,
		DueDays:    nil,
	})
	if err != nil {
		t.Fatalf("unexpected error on first borrow: %v", err)
	}

	currentLending, err := lendingQueries.GetCurrentLending(ctx, bookID)
	if err != nil {
		t.Fatalf("failed to verify precondition: %v", err)
	}
	if currentLending.BorrowerID != firstUser {
		t.Fatalf("precondition failed: expected borrower %s, got %s", firstUser, currentLending.BorrowerID)
	}

	_, err = service.BorrowBook(ctx, BorrowBookInput{
		BookID:     bookID,
		BorrowerID: secondUser,
		DueDays:    nil,
	})

	if !errors.Is(err, ErrBookAlreadyBorrowed) {
		t.Errorf("expected ErrBookAlreadyBorrowed, got %v", err)
	}

	stillBorrowed, err := lendingQueries.GetCurrentLending(ctx, bookID)
	if err != nil {
		t.Fatalf("postcondition failed: book should still be borrowed, got error: %v", err)
	}
	if stillBorrowed.LendingID != firstBorrow.ID {
		t.Errorf("postcondition failed: expected lending ID %s to remain, got %s", firstBorrow.ID, stillBorrowed.LendingID)
	}
	if stillBorrowed.BorrowerID != firstUser {
		t.Errorf("postcondition failed: expected borrower %s to remain, got %s", firstUser, stillBorrowed.BorrowerID)
	}
}

// When BorrowBook with already borrowed book by same user then extends due date
func TestBorrowBook_WithAlreadyBorrowedBookBySameUser_ExtendsDueDate(t *testing.T) {
	db := setupTestDB(t)
	lendingQueries := New(db)
	bookID := uuid.New().String()
	userID := uuid.New().String()
	bookQueries := &fakeBookQueries{
		countByBookId: map[string]int64{bookID: 1},
	}
	ctx := context.Background()

	borrowTime := time.Now().Add(-time.Duration(rand.Intn(720)+1) * time.Hour)
	service := NewBorrowBookService(lendingQueries, bookQueries)
	service.now = func() time.Time { return borrowTime }

	output1, err := service.BorrowBook(ctx, BorrowBookInput{
		BookID:     bookID,
		BorrowerID: userID,
		DueDays:    nil,
	})
	if err != nil {
		t.Fatalf("unexpected error on first borrow: %v", err)
	}

	output2, err := service.BorrowBook(ctx, BorrowBookInput{
		BookID:     bookID,
		BorrowerID: userID,
		DueDays:    nil,
	})
	if err != nil {
		t.Fatalf("unexpected error on extension: %v", err)
	}

	if output2.ID != output1.ID {
		t.Errorf("expected lending ID to remain %s, got %s", output1.ID, output2.ID)
	}

	expectedNewDueDate := output1.DueDate.AddDate(0, 0, 7).Truncate(time.Second)
	if !output2.DueDate.Truncate(time.Second).Equal(expectedNewDueDate) {
		t.Errorf("expected DueDate %v, got %v", expectedNewDueDate, output2.DueDate.Truncate(time.Second))
	}
}

// When BorrowBook with invalid dueDays then returns error
func TestBorrowBook_WithInvalidDueDays_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	lendingQueries := New(db)
	bookID := uuid.New().String()
	userID := uuid.New().String()
	bookQueries := &fakeBookQueries{
		countByBookId: map[string]int64{bookID: 1},
	}
	ctx := context.Background()

	service := NewBorrowBookService(lendingQueries, bookQueries)

	dueDays := 0
	input := BorrowBookInput{
		BookID:     bookID,
		BorrowerID: userID,
		DueDays:    &dueDays,
	}

	_, err := service.BorrowBook(ctx, input)

	if err == nil {
		t.Error("expected error, got nil")
	}
}
