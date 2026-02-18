package books

import (
	"context"
	"database/sql"

	"holocron/internal/books/domain"
)

func FindByCodeSource(queries *Queries, code *string) domain.BookListSource {
	return func(ctx context.Context, keyword *domain.SearchKeyword, pagination domain.Pagination) ([]domain.BookItem, int64, error) {
		if code == nil {
			return nil, 0, domain.ErrNotMyResponsibility
		}

		codeParam := sql.NullString{String: *code, Valid: true}

		rows, err := queries.FindBooksByCode(ctx, FindBooksByCodeParams{
			Code:   codeParam,
			Limit:  int64(pagination.Limit()),
			Offset: int64(pagination.Offset()),
		})
		if err != nil {
			return nil, 0, err
		}

		total, err := queries.CountBooksByCode(ctx, codeParam)
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
				row.CreatedAt,
				nullStringToPtr(row.BorrowerID),
				nullStringToPtr(row.BorrowerName),
				nullStringToPtr(row.BorrowedAt),
			)
			if err != nil {
				return nil, 0, err
			}
			items = append(items, *item)
		}
		return items, total, nil
	}
}
