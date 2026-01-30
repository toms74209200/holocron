package books

import (
	"context"

	"holocron/internal/books/domain"
)

func ListAllBooksSource(queries *Queries) domain.BookListSource {
	return func(ctx context.Context, keyword *domain.SearchKeyword, pagination domain.Pagination) ([]domain.BookItem, int64, error) {
		rows, err := queries.ListBooks(ctx, ListBooksParams{
			Limit:  int64(pagination.Limit()),
			Offset: int64(pagination.Offset()),
		})
		if err != nil {
			return nil, 0, err
		}

		total, err := queries.CountBooks(ctx)
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

