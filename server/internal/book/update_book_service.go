package book

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"holocron/internal/book/domain"

	"github.com/google/uuid"
)


type UpdateBookInput struct {
	BookID        string
	Code          *string
	Title         *string
	Authors       *[]string
	Publisher     *string
	PublishedDate *string
	ThumbnailURL  *string
}

type UpdateBookOutput struct {
	ID            string
	Code          *string
	Title         string
	Authors       []string
	Publisher     *string
	PublishedDate *string
	ThumbnailURL  *string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func UpdateBook(ctx context.Context, queries *Queries, input UpdateBookInput) (*UpdateBookOutput, error) {
	currentBook, err := queries.GetBookByBookId(ctx, input.BookID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBookNotFound
		}
		return nil, err
	}

	var currentAuthors []string
	if currentBook.Authors.Valid {
		if err := json.Unmarshal([]byte(currentBook.Authors.String), &currentAuthors); err != nil {
			return nil, ErrInvalidBookRow
		}
	}

	updatedCode := currentBook.Code
	if input.Code != nil {
		updatedCode = sql.NullString{String: *input.Code, Valid: true}
	}

	updatedTitle := currentBook.Title
	if input.Title != nil {
		parsedTitle, err := domain.ParseBookTitle(*input.Title)
		if err != nil {
			return nil, err
		}
		updatedTitle = sql.NullString{String: string(parsedTitle), Valid: true}
	}

	updatedAuthors := currentAuthors
	if input.Authors != nil {
		parsedAuthors, err := domain.ParseBookAuthors(*input.Authors)
		if err != nil {
			return nil, err
		}
		updatedAuthors = parsedAuthors
	}

	updatedPublisher := currentBook.Publisher
	if input.Publisher != nil {
		updatedPublisher = sql.NullString{String: *input.Publisher, Valid: true}
	}

	updatedPublishedDate := currentBook.PublishedDate
	if input.PublishedDate != nil {
		updatedPublishedDate = sql.NullString{String: *input.PublishedDate, Valid: true}
	}

	updatedThumbnailURL := currentBook.ThumbnailUrl
	if input.ThumbnailURL != nil {
		updatedThumbnailURL = sql.NullString{String: *input.ThumbnailURL, Valid: true}
	}

	authorsJSON, err := json.Marshal(updatedAuthors)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	eventID := uuid.New().String()

	err = queries.InsertBookUpdateEvent(ctx, InsertBookUpdateEventParams{
		EventID:       eventID,
		BookID:        input.BookID,
		Code:          updatedCode,
		Title:         updatedTitle,
		Authors:       sql.NullString{String: string(authorsJSON), Valid: true},
		Publisher:     updatedPublisher,
		PublishedDate: updatedPublishedDate,
		ThumbnailUrl:  updatedThumbnailURL,
		OccurredAt:    now.Format(time.RFC3339),
	})
	if err != nil {
		return nil, err
	}

	createdAtStr, ok := currentBook.CreatedAt.(string)
	if !ok || createdAtStr == "" {
		return nil, ErrInvalidBookRow
	}
	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, ErrInvalidBookRow
	}

	output := &UpdateBookOutput{
		ID:        input.BookID,
		Title:     updatedTitle.String,
		Authors:   updatedAuthors,
		Status:    "available",
		CreatedAt: createdAt,
		UpdatedAt: now,
	}

	if updatedCode.Valid {
		output.Code = &updatedCode.String
	}
	if updatedPublisher.Valid {
		output.Publisher = &updatedPublisher.String
	}
	if updatedPublishedDate.Valid {
		output.PublishedDate = &updatedPublishedDate.String
	}
	if updatedThumbnailURL.Valid {
		output.ThumbnailURL = &updatedThumbnailURL.String
	}

	return output, nil
}
