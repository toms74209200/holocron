package user

import (
	"context"
	"database/sql"

	"holocron/internal/lending"
	userDomain "holocron/internal/user/domain"
)

type GetMyBorrowingInput struct {
	BorrowerID string
}

type GetMyBorrowingOutput struct {
	Items []userDomain.BorrowingBook
	Total int64
}

func GetMyBorrowing(ctx context.Context, lendingQueries *lending.Queries, input GetMyBorrowingInput) (*GetMyBorrowingOutput, error) {
	rows, err := lendingQueries.ListBorrowingBooksByBorrowerID(ctx, input.BorrowerID)
	if err != nil {
		return nil, err
	}

	items := make([]userDomain.BorrowingBook, 0, len(rows))
	for _, row := range rows {
		book, err := userDomain.ParseBorrowingBook(
			row.BookID,
			nullStringVal(row.Title),
			nullStringPtr(row.Authors),
			row.BorrowedAt,
			nullStringPtr(row.Code),
			nullStringPtr(row.Publisher),
			nullStringPtr(row.PublishedDate),
			nullStringPtr(row.ThumbnailUrl),
			nullStringPtr(row.DueDate),
		)
		if err != nil {
			return nil, err
		}
		items = append(items, *book)
	}

	return &GetMyBorrowingOutput{
		Items: items,
		Total: int64(len(items)),
	}, nil
}

func nullStringVal(ns sql.NullString) string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}

func nullStringPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}
