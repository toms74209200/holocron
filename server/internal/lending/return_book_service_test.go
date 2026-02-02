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

func TestReturnBook_WithBorrowedBook_ReturnsOutput(t *testing.T) {
	db := setupTestDB(t)
	lendingQueries := New(db)
	bookID := uuid.New().String()
	userID := uuid.New().String()
	bookQueries := &fakeBookQueries{
		countByBookId: map[string]int64{bookID: 1},
	}
	ctx := context.Background()

	borrowTime := time.Now().Add(-time.Duration(rand.Intn(720)+1) * time.Hour)
	borrowService := NewBorrowBookService(lendingQueries, bookQueries)
	borrowService.now = func() time.Time { return borrowTime }

	borrowOutput, err := borrowService.BorrowBook(ctx, BorrowBookInput{
		BookID:     bookID,
		BorrowerID: userID,
		DueDays:    nil,
	})
	if err != nil {
		t.Fatalf("failed to borrow book: %v", err)
	}

	currentLending, err := lendingQueries.GetCurrentLending(ctx, bookID)
	if err != nil {
		t.Fatalf("failed to get current lending before return: %v", err)
	}
	if currentLending.BookID != bookID {
		t.Errorf("precondition failed: expected current lending for book %s, got %s", bookID, currentLending.BookID)
	}
	if currentLending.BorrowerID != userID {
		t.Errorf("precondition failed: expected borrower %s, got %s", userID, currentLending.BorrowerID)
	}

	returnTime := borrowTime.Add(time.Duration(rand.Intn(24)+1) * time.Hour)
	returnService := NewReturnBookService(lendingQueries, bookQueries)
	returnService.now = func() time.Time { return returnTime }

	output, err := returnService.ReturnBook(ctx, ReturnBookInput{
		BookID:      bookID,
		RequesterID: userID,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.LendingID != borrowOutput.ID {
		t.Errorf("expected LendingID %s, got %s", borrowOutput.ID, output.LendingID)
	}
	if output.BookID != bookID {
		t.Errorf("expected BookID %s, got %s", bookID, output.BookID)
	}
	if output.BorrowerID != userID {
		t.Errorf("expected BorrowerID %s, got %s", userID, output.BorrowerID)
	}
	if !output.ReturnedAt.Equal(returnTime) {
		t.Errorf("expected ReturnedAt %v, got %v", returnTime, output.ReturnedAt)
	}

	_, err = lendingQueries.GetCurrentLending(ctx, bookID)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("postcondition failed: expected no current lending after return, got error: %v", err)
	}
}

func TestReturnBook_WithNonExistentBook_ReturnsError(t *testing.T) {
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

	service := NewReturnBookService(lendingQueries, bookQueries)

	_, err = service.ReturnBook(ctx, ReturnBookInput{
		BookID:      nonExistentBookID,
		RequesterID: userID,
	})

	if !errors.Is(err, ErrBookNotFound) {
		t.Errorf("expected ErrBookNotFound, got %v", err)
	}
}

func TestReturnBook_WithNotBorrowedBook_ReturnsError(t *testing.T) {
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

	service := NewReturnBookService(lendingQueries, bookQueries)

	_, err = service.ReturnBook(ctx, ReturnBookInput{
		BookID:      bookID,
		RequesterID: userID,
	})

	if !errors.Is(err, ErrBookNotBorrowed) {
		t.Errorf("expected ErrBookNotBorrowed, got %v", err)
	}
}

func TestReturnBook_WithDifferentUser_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	lendingQueries := New(db)
	bookID := uuid.New().String()
	borrower := uuid.New().String()
	otherUser := uuid.New().String()
	bookQueries := &fakeBookQueries{
		countByBookId: map[string]int64{bookID: 1},
	}
	ctx := context.Background()

	borrowTime := time.Now().Add(-time.Duration(rand.Intn(720)+1) * time.Hour)
	borrowService := NewBorrowBookService(lendingQueries, bookQueries)
	borrowService.now = func() time.Time { return borrowTime }

	borrowOutput, err := borrowService.BorrowBook(ctx, BorrowBookInput{
		BookID:     bookID,
		BorrowerID: borrower,
		DueDays:    nil,
	})
	if err != nil {
		t.Fatalf("failed to borrow book: %v", err)
	}

	currentLending, err := lendingQueries.GetCurrentLending(ctx, bookID)
	if err != nil {
		t.Fatalf("failed to verify precondition: %v", err)
	}
	if currentLending.BorrowerID != borrower {
		t.Fatalf("precondition failed: expected borrower %s, got %s", borrower, currentLending.BorrowerID)
	}

	returnService := NewReturnBookService(lendingQueries, bookQueries)

	_, err = returnService.ReturnBook(ctx, ReturnBookInput{
		BookID:      bookID,
		RequesterID: otherUser,
	})

	if !errors.Is(err, ErrNotBorrower) {
		t.Errorf("expected ErrNotBorrower, got %v", err)
	}

	stillBorrowed, err := lendingQueries.GetCurrentLending(ctx, bookID)
	if err != nil {
		t.Fatalf("postcondition failed: book should still be borrowed, got error: %v", err)
	}
	if stillBorrowed.LendingID != borrowOutput.ID {
		t.Errorf("postcondition failed: expected lending ID %s to remain, got %s", borrowOutput.ID, stillBorrowed.LendingID)
	}
}
