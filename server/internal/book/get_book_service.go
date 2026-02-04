package book

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrBookNotFound   = errors.New("book not found")
	ErrInvalidBookID  = errors.New("invalid book ID")
	ErrInvalidBookRow = errors.New("invalid book row")
)

type GetBookInput struct {
	BookID string
}

type Borrower struct {
	ID         string
	Name       string
	BorrowedAt time.Time
}

type GetBookOutput struct {
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

func GetBook(ctx context.Context, queries *Queries, input GetBookInput) (*GetBookOutput, error) {
	if input.BookID == "" {
		return nil, ErrInvalidBookID
	}

	row, err := queries.GetBookByBookId(ctx, input.BookID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBookNotFound
		}
		return nil, err
	}

	var authors []string
	if row.Authors.Valid {
		if err := json.Unmarshal([]byte(row.Authors.String), &authors); err != nil {
			return nil, ErrInvalidBookRow
		}
	}

	createdAt, err := time.Parse(time.RFC3339, row.OccurredAt)
	if err != nil {
		return nil, ErrInvalidBookRow
	}

	// Check if the book is currently borrowed
	status := "available"
	var borrower *Borrower
	borrowerInfo, err := queries.GetBookBorrowerInfo(ctx, input.BookID)
	if err == nil {
		// Book is borrowed
		status = "borrowed"
		borrowedAt, parseErr := time.Parse(time.RFC3339, borrowerInfo.BorrowedAt)
		if parseErr != nil {
			return nil, ErrInvalidBookRow
		}
		borrower = &Borrower{
			ID:         borrowerInfo.BorrowerID,
			Name:       borrowerInfo.BorrowerName,
			BorrowedAt: borrowedAt,
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		// Unexpected error (not "no rows" error)
		return nil, err
	}

	return &GetBookOutput{
		ID:            row.BookID,
		Code:          nullStringToPtr(row.Code),
		Title:         row.Title.String,
		Authors:       authors,
		Publisher:     nullStringToPtr(row.Publisher),
		PublishedDate: nullStringToPtr(row.PublishedDate),
		ThumbnailURL:  nullStringToPtr(row.ThumbnailUrl),
		Status:        status,
		Borrower:      borrower,
		CreatedAt:     createdAt,
	}, nil
}

func nullStringToPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}
