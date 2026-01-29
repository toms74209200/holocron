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

var ErrInvalidCode = errors.New("invalid code")

type CreateBookByCodeInput struct {
	Code string
}

type CreateBookByCodeOutput struct {
	ID            string
	Code          string
	Title         string
	Authors       []string
	Publisher     *string
	PublishedDate *string
	ThumbnailURL  *string
	Status        string
	CreatedAt     time.Time
}

func DBCacheSource(queries *Queries) domain.BookInfoSource {
	return func(ctx context.Context, code string) (*domain.BookInfo, error) {
		row, err := queries.GetBookByCode(ctx, sql.NullString{String: code, Valid: true})
		if err != nil {
			return nil, err
		}

		var authors []string
		if row.Authors.Valid && row.Authors.String != "" {
			_ = json.Unmarshal([]byte(row.Authors.String), &authors)
		}

		return &domain.BookInfo{
			Title:         row.Title.String,
			Authors:       authors,
			Publisher:     row.Publisher.String,
			PublishedDate: row.PublishedDate.String,
			ThumbnailURL:  row.ThumbnailUrl.String,
		}, nil
	}
}

func ExternalAPISource(
	fetch func(ctx context.Context, code string) ([]byte, error),
	parse func(body []byte) (*domain.BookInfo, error),
) domain.BookInfoSource {
	return func(ctx context.Context, code string) (*domain.BookInfo, error) {
		body, err := fetch(ctx, code)
		if err != nil {
			return nil, err
		}
		return parse(body)
	}
}

func CreateBookByCode(
	ctx context.Context,
	queries *Queries,
	sources []domain.BookInfoSource,
	input CreateBookByCodeInput,
) (*CreateBookByCodeOutput, error) {
	code, err := domain.ParseBookCode(input.Code)
	if err != nil {
		return nil, ErrInvalidCode
	}

	info, err := domain.LookupBookInfo(ctx, sources, string(code))
	if err != nil {
		return nil, err
	}

	title, err := domain.ParseBookTitle(info.Title)
	if err != nil {
		return nil, ErrInvalidTitle
	}

	authors, err := domain.ParseBookAuthors(info.Authors)
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
		Code:          sql.NullString{String: string(code), Valid: true},
		Title:         sql.NullString{String: string(title), Valid: true},
		Authors:       sql.NullString{String: string(authorsJSON), Valid: true},
		Publisher:     toNullString(strPtr(info.Publisher)),
		PublishedDate: toNullString(strPtr(info.PublishedDate)),
		ThumbnailUrl:  toNullString(strPtr(info.ThumbnailURL)),
		OccurredAt:    now.Format(time.RFC3339),
	})
	if err != nil {
		return nil, err
	}

	return &CreateBookByCodeOutput{
		ID:            bookID,
		Code:          string(code),
		Title:         string(title),
		Authors:       authors,
		Publisher:     strPtr(info.Publisher),
		PublishedDate: strPtr(info.PublishedDate),
		ThumbnailURL:  strPtr(info.ThumbnailURL),
		Status:        "available",
		CreatedAt:     now,
	}, nil
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
