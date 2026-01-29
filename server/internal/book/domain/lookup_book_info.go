package domain

import (
	"context"
)

type BookInfoSource func(ctx context.Context, code string) (*BookInfo, error)

func LookupBookInfo(ctx context.Context, sources []BookInfoSource, code string) (*BookInfo, error) {
	for _, src := range sources {
		info, err := src(ctx, code)
		if err == nil {
			return info, nil
		}
	}
	return nil, ErrBookNotFound
}
