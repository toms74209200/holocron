package books

import (
	"context"
	"database/sql"

	"holocron/internal/books/domain"
)

func SearchBooksSource(queries *Queries) domain.BookListSource {
	return func(ctx context.Context, keyword *domain.SearchKeyword, pagination domain.Pagination) ([]domain.BookItem, int64, error) {
		if keyword == nil {
			return nil, 0, domain.ErrNotMyResponsibility
		}

		pattern := "%" + string(*keyword) + "%"
		searchPattern := sql.NullString{String: pattern, Valid: true}

		rows, err := queries.SearchBooks(ctx, SearchBooksParams{
			Title:   searchPattern,
			Authors: searchPattern,
			Limit:   int64(pagination.Limit()),
			Offset:  int64(pagination.Offset()),
		})
		if err != nil {
			return nil, 0, err
		}

		total, err := queries.CountSearchBooks(ctx, CountSearchBooksParams{
			Title:   searchPattern,
			Authors: searchPattern,
		})
		if err != nil {
			return nil, 0, err
		}

		items := make([]domain.BookItem, 0, len(rows))
		for _, row := range rows {
			item, err := domain.BookItemFromRow(
				row.BookID,
				nullStringToPtr(row.Code),
				row.Title.String,
				row.Authors.String,
				nullStringToPtr(row.Publisher),
				nullStringToPtr(row.PublishedDate),
				nullStringToPtr(row.ThumbnailUrl),
				row.OccurredAt,
			)
			if err != nil {
				return nil, 0, err
			}
			items = append(items, *item)
		}
		return items, total, nil
	}
}

func nullStringToPtr(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}
