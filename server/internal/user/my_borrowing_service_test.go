//go:build medium

package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"

	"holocron/internal/lending"
)

func setupBorrowingTestDB(t *testing.T) *sql.DB {
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
			delete_reason TEXT,
			delete_memo TEXT,
			occurred_at TEXT NOT NULL
		);
		CREATE INDEX idx_book_events_book_id ON book_events(book_id);
	`)
	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func insertTestBookEvent(t *testing.T, db *sql.DB, bookID, title string, authors []string) {
	t.Helper()
	authorsJSON, err := json.Marshal(authors)
	if err != nil {
		t.Fatalf("failed to marshal authors: %v", err)
	}
	_, err = db.Exec(
		`INSERT INTO book_events (event_id, book_id, event_type, title, authors, occurred_at) VALUES (?, ?, 'created', ?, ?, ?)`,
		uuid.New().String(), bookID, title, string(authorsJSON), time.Now().UTC().Format(time.RFC3339),
	)
	if err != nil {
		t.Fatalf("failed to insert book event: %v", err)
	}
}

func insertTestBorrowEvent(t *testing.T, db *sql.DB, bookID, borrowerID, lendingID string, borrowedAt, dueDate time.Time) {
	t.Helper()
	_, err := db.Exec(
		`INSERT INTO lending_events (event_id, lending_id, book_id, borrower_id, event_type, due_date, occurred_at) VALUES (?, ?, ?, ?, 'borrowed', ?, ?)`,
		uuid.New().String(), lendingID, bookID, borrowerID,
		dueDate.UTC().Format(time.RFC3339), borrowedAt.UTC().Format(time.RFC3339),
	)
	if err != nil {
		t.Fatalf("failed to insert borrow event: %v", err)
	}
}

func insertTestReturnEvent(t *testing.T, db *sql.DB, bookID, borrowerID, lendingID string, returnedAt time.Time) {
	t.Helper()
	_, err := db.Exec(
		`INSERT INTO lending_events (event_id, lending_id, book_id, borrower_id, event_type, due_date, occurred_at) VALUES (?, ?, ?, ?, 'returned', NULL, ?)`,
		uuid.New().String(), lendingID, bookID, borrowerID, returnedAt.UTC().Format(time.RFC3339),
	)
	if err != nil {
		t.Fatalf("failed to insert return event: %v", err)
	}
}

// When GetMyBorrowing with no borrowed books then returns empty list
func TestGetMyBorrowing_WithNoBorrowedBooks_ReturnsEmptyList(t *testing.T) {
	db := setupBorrowingTestDB(t)
	lendingQueries := lending.New(db)
	borrowerID := uuid.New().String()
	ctx := context.Background()

	rows, err := lendingQueries.ListBorrowingBooksByBorrowerID(ctx, borrowerID)
	if err != nil {
		t.Fatalf("precondition failed: %v", err)
	}
	if len(rows) != 0 {
		t.Fatalf("precondition failed: expected no rows, got %d", len(rows))
	}

	output, err := GetMyBorrowing(ctx, lendingQueries, GetMyBorrowingInput{
		BorrowerID: borrowerID,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(output.Items) != 0 {
		t.Errorf("expected no items, got %d", len(output.Items))
	}
	if output.Total != 0 {
		t.Errorf("expected total 0, got %d", output.Total)
	}
}

// When GetMyBorrowing with one borrowed book then returns the book
func TestGetMyBorrowing_WithOneBorrowedBook_ReturnsBook(t *testing.T) {
	db := setupBorrowingTestDB(t)
	lendingQueries := lending.New(db)
	borrowerID := uuid.New().String()
	bookID := uuid.New().String()
	lendingID := uuid.New().String()
	ctx := context.Background()

	borrowedAt := time.Now().Add(-2 * time.Hour).UTC().Truncate(time.Second)
	dueDate := borrowedAt.AddDate(0, 0, 7)

	insertTestBookEvent(t, db, bookID, "テスト書籍", []string{"著者A"})
	insertTestBorrowEvent(t, db, bookID, borrowerID, lendingID, borrowedAt, dueDate)

	rows, err := lendingQueries.ListBorrowingBooksByBorrowerID(ctx, borrowerID)
	if err != nil {
		t.Fatalf("precondition failed: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("precondition failed: expected 1 row, got %d", len(rows))
	}

	output, err := GetMyBorrowing(ctx, lendingQueries, GetMyBorrowingInput{
		BorrowerID: borrowerID,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(output.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(output.Items))
	}
	if output.Total != 1 {
		t.Errorf("expected total 1, got %d", output.Total)
	}

	item := output.Items[0]
	if item.ID != bookID {
		t.Errorf("expected ID %s, got %s", bookID, item.ID)
	}
	if item.Title != "テスト書籍" {
		t.Errorf("expected Title テスト書籍, got %s", item.Title)
	}
	if len(item.Authors) != 1 || item.Authors[0] != "著者A" {
		t.Errorf("expected Authors [著者A], got %v", item.Authors)
	}
	if !item.BorrowedAt.Equal(borrowedAt) {
		t.Errorf("expected BorrowedAt %v, got %v", borrowedAt, item.BorrowedAt)
	}
	if item.DueDate == nil {
		t.Error("expected DueDate to be set")
	} else if !item.DueDate.Equal(dueDate) {
		t.Errorf("expected DueDate %v, got %v", dueDate, *item.DueDate)
	}
}

// When GetMyBorrowing after returning book then returns empty list
func TestGetMyBorrowing_AfterReturningBook_ReturnsEmptyList(t *testing.T) {
	db := setupBorrowingTestDB(t)
	lendingQueries := lending.New(db)
	borrowerID := uuid.New().String()
	bookID := uuid.New().String()
	lendingID := uuid.New().String()
	ctx := context.Background()

	borrowedAt := time.Now().Add(-2 * time.Hour).UTC().Truncate(time.Second)
	dueDate := borrowedAt.AddDate(0, 0, 7)
	returnedAt := borrowedAt.Add(1 * time.Hour)

	insertTestBookEvent(t, db, bookID, "テスト書籍", []string{"著者A"})
	insertTestBorrowEvent(t, db, bookID, borrowerID, lendingID, borrowedAt, dueDate)
	insertTestReturnEvent(t, db, bookID, borrowerID, lendingID, returnedAt)

	output, err := GetMyBorrowing(ctx, lendingQueries, GetMyBorrowingInput{
		BorrowerID: borrowerID,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(output.Items) != 0 {
		t.Errorf("expected no items after return, got %d", len(output.Items))
	}
	if output.Total != 0 {
		t.Errorf("expected total 0, got %d", output.Total)
	}
}

// When GetMyBorrowing with other user's borrowed book then returns empty list
func TestGetMyBorrowing_WithOtherUserBorrowingBook_ReturnsEmptyList(t *testing.T) {
	db := setupBorrowingTestDB(t)
	lendingQueries := lending.New(db)
	myUserID := uuid.New().String()
	otherUserID := uuid.New().String()
	bookID := uuid.New().String()
	lendingID := uuid.New().String()
	ctx := context.Background()

	borrowedAt := time.Now().Add(-2 * time.Hour).UTC().Truncate(time.Second)
	dueDate := borrowedAt.AddDate(0, 0, 7)

	insertTestBookEvent(t, db, bookID, "他人の書籍", []string{"著者B"})
	insertTestBorrowEvent(t, db, bookID, otherUserID, lendingID, borrowedAt, dueDate)

	output, err := GetMyBorrowing(ctx, lendingQueries, GetMyBorrowingInput{
		BorrowerID: myUserID,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(output.Items) != 0 {
		t.Errorf("expected no items for my user, got %d", len(output.Items))
	}
}
