package book

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"holocron/internal/book/domain"

	"github.com/google/uuid"
)

var (
	ErrBookBorrowed        = errors.New("book is currently borrowed")
	ErrInvalidDeleteReason = errors.New("invalid delete reason")
)

type DeleteBookInput struct {
	BookID string
	Reason string
	Memo   *string
}

func DeleteBook(ctx context.Context, queries *Queries, input DeleteBookInput) error {
	reason, err := domain.ParseDeleteReason(input.Reason)
	if err != nil {
		return ErrInvalidDeleteReason
	}

	_, err = queries.GetBookByBookId(ctx, input.BookID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrBookNotFound
		}
		return err
	}

	_, err = queries.GetBookBorrowerInfo(ctx, input.BookID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if err == nil {
		return ErrBookBorrowed
	}

	now := time.Now().UTC()
	eventID := uuid.New().String()

	var memo sql.NullString
	if input.Memo != nil {
		memo = sql.NullString{String: *input.Memo, Valid: true}
	}

	err = queries.InsertBookDeleteEvent(ctx, InsertBookDeleteEventParams{
		EventID:      eventID,
		BookID:       input.BookID,
		DeleteReason: sql.NullString{String: string(reason), Valid: true},
		DeleteMemo:   memo,
		OccurredAt:   now.Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	return nil
}
