package domain

import (
	"context"
	"errors"
)

var ErrBookListNotAvailable = errors.New("book list not available")
var ErrNotMyResponsibility = errors.New("not my responsibility")

type BookListSource func(ctx context.Context, keyword *SearchKeyword, pagination Pagination) ([]BookItem, int64, error)

func GetBookList(ctx context.Context, sources []BookListSource, keyword *SearchKeyword, pagination Pagination) ([]BookItem, int64, error) {
	for _, src := range sources {
		items, total, err := src(ctx, keyword, pagination)
		if err == nil {
			return items, total, nil
		}
	}
	return nil, 0, ErrBookListNotAvailable
}
