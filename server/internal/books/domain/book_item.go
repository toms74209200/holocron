package domain

import (
	"encoding/json"
	"errors"
	"time"
)

var ErrInvalidBookRow = errors.New("invalid book row")

type Borrower struct {
	ID         string
	Name       string
	BorrowedAt time.Time
}

type BookItem struct {
	ID            string
	Code          *string
	Title         string
	Authors       []string
	Publisher     *string
	PublishedDate *string
	ThumbnailURL  *string
	Status        string
	Borrower      *Borrower
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
	createdAtRaw interface{},
	borrowerID *string,
	borrowerName *string,
	borrowedAt *string,
) (*BookItem, error) {
	var authors []string
	if authorsJSON != "" {
		if err := json.Unmarshal([]byte(authorsJSON), &authors); err != nil {
			return nil, ErrInvalidBookRow
		}
	}

	createdAtStr, ok := createdAtRaw.(string)
	if !ok || createdAtStr == "" {
		return nil, ErrInvalidBookRow
	}
	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, ErrInvalidBookRow
	}

	status := "available"
	var borrower *Borrower
	if borrowerID != nil && borrowerName != nil && borrowedAt != nil {
		status = "borrowed"
		borrowedAtTime, parseErr := time.Parse(time.RFC3339, *borrowedAt)
		if parseErr != nil {
			return nil, ErrInvalidBookRow
		}
		borrower = &Borrower{
			ID:         *borrowerID,
			Name:       *borrowerName,
			BorrowedAt: borrowedAtTime,
		}
	}

	return &BookItem{
		ID:            bookID,
		Code:          code,
		Title:         title,
		Authors:       authors,
		Publisher:     publisher,
		PublishedDate: publishedDate,
		ThumbnailURL:  thumbnailURL,
		Status:        status,
		Borrower:      borrower,
		CreatedAt:     createdAt,
	}, nil
}
