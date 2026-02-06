//go:build medium

package book

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func setupDeleteTestDB(t *testing.T) *sql.DB {
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
			delete_reason TEXT,
			delete_memo TEXT,
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

		CREATE TABLE user_events (
			event_id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			event_type TEXT NOT NULL,
			name TEXT NOT NULL,
			occurred_at TEXT NOT NULL
		);
		CREATE INDEX idx_user_events_user_id ON user_events(user_id);
	`)
	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

// When DeleteBook with available book then deletes successfully
func TestDeleteBook_WithAvailableBook_DeletesSuccessfully(t *testing.T) {
	db := setupDeleteTestDB(t)
	queries := New(db)
	ctx := context.Background()

	bookID := uuid.New().String()
	title := uuid.New().String()
	authors := fmt.Sprintf(`["%s"]`, uuid.New().String())
	occurredAt := time.Now().Add(-time.Duration(rand.Intn(720)+1) * time.Hour).Format(time.RFC3339)
	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, title, authors, occurred_at)
		VALUES (?, ?, 'created', ?, ?, ?)
	`, uuid.New().String(), bookID, title, authors, occurredAt)
	if err != nil {
		t.Fatalf("failed to insert book: %v", err)
	}

	err = DeleteBook(ctx, queries, DeleteBookInput{
		BookID: bookID,
		Reason: "disposal",
		Memo:   nil,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var eventCount int
	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM book_events WHERE book_id = ? AND event_type = 'deleted'`, bookID).Scan(&eventCount)
	if err != nil {
		t.Fatalf("failed to count events: %v", err)
	}
	expectedEventCount := 1
	if eventCount != expectedEventCount {
		t.Errorf("expected %d delete event, got %d", expectedEventCount, eventCount)
	}

	var deleteReason string
	err = db.QueryRowContext(ctx, `SELECT delete_reason FROM book_events WHERE book_id = ? AND event_type = 'deleted'`, bookID).Scan(&deleteReason)
	if err != nil {
		t.Fatalf("failed to get delete reason: %v", err)
	}
	if deleteReason != "disposal" {
		t.Errorf("expected delete reason 'disposal', got %s", deleteReason)
	}

	_, err = queries.GetBookByBookId(ctx, bookID)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected book to be deleted (sql.ErrNoRows), got %v", err)
	}
}

// When DeleteBook with memo then stores memo
func TestDeleteBook_WithMemo_StoresMemo(t *testing.T) {
	db := setupDeleteTestDB(t)
	queries := New(db)
	ctx := context.Background()

	bookID := uuid.New().String()
	title := uuid.New().String()
	authors := fmt.Sprintf(`["%s"]`, uuid.New().String())
	occurredAt := time.Now().Add(-time.Duration(rand.Intn(720)+1) * time.Hour).Format(time.RFC3339)
	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, title, authors, occurred_at)
		VALUES (?, ?, 'created', ?, ?, ?)
	`, uuid.New().String(), bookID, title, authors, occurredAt)
	if err != nil {
		t.Fatalf("failed to insert book: %v", err)
	}

	memo := "友人に譲りました"
	err = DeleteBook(ctx, queries, DeleteBookInput{
		BookID: bookID,
		Reason: "transfer",
		Memo:   &memo,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var deleteReason, deleteMemo string
	err = db.QueryRowContext(ctx, `SELECT delete_reason, delete_memo FROM book_events WHERE book_id = ? AND event_type = 'deleted'`, bookID).Scan(&deleteReason, &deleteMemo)
	if err != nil {
		t.Fatalf("failed to get delete reason and memo: %v", err)
	}
	if deleteReason != "transfer" {
		t.Errorf("expected delete reason 'transfer', got %s", deleteReason)
	}
	if deleteMemo != memo {
		t.Errorf("expected delete memo '%s', got %s", memo, deleteMemo)
	}
}

// When DeleteBook with non-existent book then returns not found error
func TestDeleteBook_WithNonExistentBook_ReturnsNotFoundError(t *testing.T) {
	db := setupDeleteTestDB(t)
	queries := New(db)
	ctx := context.Background()

	nonExistentBookID := uuid.New().String()
	err := DeleteBook(ctx, queries, DeleteBookInput{
		BookID: nonExistentBookID,
		Reason: "disposal",
		Memo:   nil,
	})

	if !errors.Is(err, ErrBookNotFound) {
		t.Errorf("expected ErrBookNotFound, got %v", err)
	}
}

// When DeleteBook with borrowed book then returns book borrowed error
func TestDeleteBook_WithBorrowedBook_ReturnsBookBorrowedError(t *testing.T) {
	db := setupDeleteTestDB(t)
	queries := New(db)
	ctx := context.Background()

	bookID := uuid.New().String()
	title := uuid.New().String()
	authors := fmt.Sprintf(`["%s"]`, uuid.New().String())
	occurredAt := time.Now().Add(-time.Duration(rand.Intn(720)+1) * time.Hour).Format(time.RFC3339)
	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, title, authors, occurred_at)
		VALUES (?, ?, 'created', ?, ?, ?)
	`, uuid.New().String(), bookID, title, authors, occurredAt)
	if err != nil {
		t.Fatalf("failed to insert book: %v", err)
	}

	userID := uuid.New().String()
	userName := uuid.New().String()
	_, err = db.ExecContext(ctx, `
		INSERT INTO user_events (event_id, user_id, event_type, name, occurred_at)
		VALUES (?, ?, 'created', ?, ?)
	`, uuid.New().String(), userID, userName, occurredAt)
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	lendingID := uuid.New().String()
	borrowedAt := time.Now().Format(time.RFC3339)
	_, err = db.ExecContext(ctx, `
		INSERT INTO lending_events (event_id, lending_id, book_id, borrower_id, event_type, occurred_at)
		VALUES (?, ?, ?, ?, 'borrowed', ?)
	`, uuid.New().String(), lendingID, bookID, userID, borrowedAt)
	if err != nil {
		t.Fatalf("failed to insert lending event: %v", err)
	}

	err = DeleteBook(ctx, queries, DeleteBookInput{
		BookID: bookID,
		Reason: "disposal",
		Memo:   nil,
	})

	if !errors.Is(err, ErrBookBorrowed) {
		t.Errorf("expected ErrBookBorrowed, got %v", err)
	}

	var eventCount int
	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM book_events WHERE book_id = ? AND event_type = 'deleted'`, bookID).Scan(&eventCount)
	if err != nil {
		t.Fatalf("failed to count events: %v", err)
	}
	expectedEventCount := 0
	if eventCount != expectedEventCount {
		t.Errorf("expected %d delete events (should not create delete event for borrowed book), got %d", expectedEventCount, eventCount)
	}
}

// When DeleteBook with returned book then deletes successfully
func TestDeleteBook_WithReturnedBook_DeletesSuccessfully(t *testing.T) {
	db := setupDeleteTestDB(t)
	queries := New(db)
	ctx := context.Background()

	bookID := uuid.New().String()
	title := uuid.New().String()
	authors := fmt.Sprintf(`["%s"]`, uuid.New().String())
	occurredAt := time.Now().Add(-time.Duration(rand.Intn(720)+1) * time.Hour).Format(time.RFC3339)
	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, title, authors, occurred_at)
		VALUES (?, ?, 'created', ?, ?, ?)
	`, uuid.New().String(), bookID, title, authors, occurredAt)
	if err != nil {
		t.Fatalf("failed to insert book: %v", err)
	}

	userID := uuid.New().String()
	userName := uuid.New().String()
	_, err = db.ExecContext(ctx, `
		INSERT INTO user_events (event_id, user_id, event_type, name, occurred_at)
		VALUES (?, ?, 'created', ?, ?)
	`, uuid.New().String(), userID, userName, occurredAt)
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	lendingID := uuid.New().String()
	borrowedAt := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	_, err = db.ExecContext(ctx, `
		INSERT INTO lending_events (event_id, lending_id, book_id, borrower_id, event_type, occurred_at)
		VALUES (?, ?, ?, ?, 'borrowed', ?)
	`, uuid.New().String(), lendingID, bookID, userID, borrowedAt)
	if err != nil {
		t.Fatalf("failed to insert lending event: %v", err)
	}

	returnedAt := time.Now().Format(time.RFC3339)
	_, err = db.ExecContext(ctx, `
		INSERT INTO lending_events (event_id, lending_id, book_id, borrower_id, event_type, occurred_at)
		VALUES (?, ?, ?, ?, 'returned', ?)
	`, uuid.New().String(), lendingID, bookID, userID, returnedAt)
	if err != nil {
		t.Fatalf("failed to insert return event: %v", err)
	}

	err = DeleteBook(ctx, queries, DeleteBookInput{
		BookID: bookID,
		Reason: "lost",
		Memo:   nil,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var eventCount int
	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM book_events WHERE book_id = ? AND event_type = 'deleted'`, bookID).Scan(&eventCount)
	if err != nil {
		t.Fatalf("failed to count events: %v", err)
	}
	expectedEventCount := 1
	if eventCount != expectedEventCount {
		t.Errorf("expected %d delete event, got %d", expectedEventCount, eventCount)
	}
}

// When DeleteBook with empty reason then returns invalid delete reason error
func TestDeleteBook_WithEmptyReason_ReturnsInvalidDeleteReasonError(t *testing.T) {
	db := setupDeleteTestDB(t)
	queries := New(db)
	ctx := context.Background()

	bookID := uuid.New().String()
	title := uuid.New().String()
	authors := fmt.Sprintf(`["%s"]`, uuid.New().String())
	occurredAt := time.Now().Add(-time.Duration(rand.Intn(720)+1) * time.Hour).Format(time.RFC3339)
	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, title, authors, occurred_at)
		VALUES (?, ?, 'created', ?, ?, ?)
	`, uuid.New().String(), bookID, title, authors, occurredAt)
	if err != nil {
		t.Fatalf("failed to insert book: %v", err)
	}

	err = DeleteBook(ctx, queries, DeleteBookInput{
		BookID: bookID,
		Reason: "",
		Memo:   nil,
	})

	if !errors.Is(err, ErrInvalidDeleteReason) {
		t.Errorf("expected ErrInvalidDeleteReason, got %v", err)
	}
}

// When DeleteBook with invalid reason then returns invalid delete reason error
func TestDeleteBook_WithInvalidReason_ReturnsInvalidDeleteReasonError(t *testing.T) {
	db := setupDeleteTestDB(t)
	queries := New(db)
	ctx := context.Background()

	bookID := uuid.New().String()
	title := uuid.New().String()
	authors := fmt.Sprintf(`["%s"]`, uuid.New().String())
	occurredAt := time.Now().Add(-time.Duration(rand.Intn(720)+1) * time.Hour).Format(time.RFC3339)
	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, title, authors, occurred_at)
		VALUES (?, ?, 'created', ?, ?, ?)
	`, uuid.New().String(), bookID, title, authors, occurredAt)
	if err != nil {
		t.Fatalf("failed to insert book: %v", err)
	}

	err = DeleteBook(ctx, queries, DeleteBookInput{
		BookID: bookID,
		Reason: "invalid_reason",
		Memo:   nil,
	})

	if !errors.Is(err, ErrInvalidDeleteReason) {
		t.Errorf("expected ErrInvalidDeleteReason, got %v", err)
	}
}

// When DeleteBook with already deleted book then returns not found error
func TestDeleteBook_WithAlreadyDeletedBook_ReturnsNotFoundError(t *testing.T) {
	db := setupDeleteTestDB(t)
	queries := New(db)
	ctx := context.Background()

	bookID := uuid.New().String()
	title := uuid.New().String()
	authors := fmt.Sprintf(`["%s"]`, uuid.New().String())
	occurredAt := time.Now().Add(-time.Duration(rand.Intn(720)+1) * time.Hour).Format(time.RFC3339)
	_, err := db.ExecContext(ctx, `
		INSERT INTO book_events (event_id, book_id, event_type, title, authors, occurred_at)
		VALUES (?, ?, 'created', ?, ?, ?)
	`, uuid.New().String(), bookID, title, authors, occurredAt)
	if err != nil {
		t.Fatalf("failed to insert book: %v", err)
	}

	err = DeleteBook(ctx, queries, DeleteBookInput{
		BookID: bookID,
		Reason: "disposal",
		Memo:   nil,
	})
	if err != nil {
		t.Fatalf("first deletion failed: %v", err)
	}

	err = DeleteBook(ctx, queries, DeleteBookInput{
		BookID: bookID,
		Reason: "disposal",
		Memo:   nil,
	})

	if !errors.Is(err, ErrBookNotFound) {
		t.Errorf("expected ErrBookNotFound when deleting already deleted book, got %v", err)
	}

	var eventCount int
	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM book_events WHERE book_id = ? AND event_type = 'deleted'`, bookID).Scan(&eventCount)
	if err != nil {
		t.Fatalf("failed to count events: %v", err)
	}
	expectedEventCount := 1
	if eventCount != expectedEventCount {
		t.Errorf("expected %d delete event, got %d", expectedEventCount, eventCount)
	}
}
