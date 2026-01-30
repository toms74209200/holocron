package domain

import (
	"encoding/json"
	"errors"
	"time"
)

var ErrInvalidBookRow = errors.New("invalid book row")

type BookItem struct {
	ID            string
	Code          *string
	Title         string
	Authors       []string
	Publisher     *string
	PublishedDate *string
	ThumbnailURL  *string
	Status        string
	CreatedAt     time.Time
}

func BookItemFromRow(
	bookID string,
	code *string,
	title string,
	authorsJSON string,
	publisher *string,
	publishedDate *string,
	thumbnailURL *string,
	occurredAt string,
) (*BookItem, error) {
	var authors []string
	if authorsJSON != "" {
		if err := json.Unmarshal([]byte(authorsJSON), &authors); err != nil {
			return nil, ErrInvalidBookRow
		}
	}

	createdAt, err := time.Parse(time.RFC3339, occurredAt)
	if err != nil {
		return nil, ErrInvalidBookRow
	}

	return &BookItem{
		ID:            bookID,
		Code:          code,
		Title:         title,
		Authors:       authors,
		Publisher:     publisher,
		PublishedDate: publishedDate,
		ThumbnailURL:  thumbnailURL,
		Status:        "available",
		CreatedAt:     createdAt,
	}, nil
}
