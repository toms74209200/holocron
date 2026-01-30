package domain

import (
	"context"

	book "holocron/internal/book/domain"
)

type BookInfoSource func(ctx context.Context, code string) (*book.BookInfo, error)

func LookupBookInfo(ctx context.Context, sources []BookInfoSource, code string) (*book.BookInfo, error) {
	for _, src := range sources {
		info, err := src(ctx, code)
		if err == nil {
			return info, nil
		}
	}
	return nil, book.ErrBookNotFound
}
