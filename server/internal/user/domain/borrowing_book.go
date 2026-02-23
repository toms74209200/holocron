package domain

import (
	"encoding/json"
	"errors"
	"time"
)

var ErrInvalidBorrowingBook = errors.New("invalid borrowing book data")

type BorrowingBook struct {
	ID            string
	Code          *string
	Title         string
	Authors       []string
	Publisher     *string
	PublishedDate *string
	ThumbnailURL  *string
	BorrowedAt    time.Time
	DueDate       *time.Time
}

func ParseBorrowingBook(
	id string,
	title string,
	authorsJSON *string,
	borrowedAt string,
	code, publisher, publishedDate, thumbnailURL, dueDate *string,
) (*BorrowingBook, error) {
	var authors []string
	if authorsJSON != nil && *authorsJSON != "" {
		if err := json.Unmarshal([]byte(*authorsJSON), &authors); err != nil {
			return nil, ErrInvalidBorrowingBook
		}
	} else {
		authors = []string{}
	}

	parsedBorrowedAt, err := time.Parse(time.RFC3339, borrowedAt)
	if err != nil {
		return nil, ErrInvalidBorrowingBook
	}

	var parsedDueDate *time.Time
	if dueDate != nil && *dueDate != "" {
		t, err := time.Parse(time.RFC3339, *dueDate)
		if err != nil {
			return nil, ErrInvalidBorrowingBook
		}
		parsedDueDate = &t
	}

	return &BorrowingBook{
		ID:            id,
		Code:          code,
		Title:         title,
		Authors:       authors,
		Publisher:     publisher,
		PublishedDate: publishedDate,
		ThumbnailURL:  thumbnailURL,
		BorrowedAt:    parsedBorrowedAt,
		DueDate:       parsedDueDate,
	}, nil
}
