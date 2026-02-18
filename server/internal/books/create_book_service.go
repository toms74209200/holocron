package books

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	book "holocron/internal/book/domain"

	"github.com/google/uuid"
)

var (
	ErrInvalidTitle   = errors.New("invalid title")
	ErrInvalidAuthors = errors.New("invalid authors")
)

type CreateBookInput struct {
	Code          *string
	Title         string
	Authors       []string
	Publisher     *string
	PublishedDate *string
	ThumbnailURL  *string
}

type CreateBookOutput struct {
	ID            string
	Title         string
	Authors       []string
	Publisher     *string
	PublishedDate *string
	ThumbnailURL  *string
	Status        string
	CreatedAt     time.Time
}

func CreateBook(ctx context.Context, queries *Queries, input CreateBookInput) (*CreateBookOutput, error) {
	title, err := book.ParseBookTitle(input.Title)
	if err != nil {
		return nil, ErrInvalidTitle
	}

	authors, err := book.ParseBookAuthors(input.Authors)
	if err != nil {
		return nil, ErrInvalidAuthors
	}

	bookID := uuid.New().String()
	now := time.Now().UTC()

	authorsJSON, err := json.Marshal(authors)
	if err != nil {
		return nil, err
	}

	err = queries.InsertBookEvent(ctx, InsertBookEventParams{
		EventID:       uuid.New().String(),
		BookID:        bookID,
		EventType:     "created",
		Code:          toNullString(input.Code),
		Title:         sql.NullString{String: string(title), Valid: true},
		Authors:       sql.NullString{String: string(authorsJSON), Valid: true},
		Publisher:     toNullString(input.Publisher),
		PublishedDate: toNullString(input.PublishedDate),
		ThumbnailUrl:  toNullString(input.ThumbnailURL),
		OccurredAt:    now.Format(time.RFC3339),
	})
	if err != nil {
		return nil, err
	}

	return &CreateBookOutput{
		ID:            bookID,
		Title:         string(title),
		Authors:       authors,
		Publisher:     input.Publisher,
		PublishedDate: input.PublishedDate,
		ThumbnailURL:  input.ThumbnailURL,
		Status:        "available",
		CreatedAt:     now,
	}, nil
}

func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}
