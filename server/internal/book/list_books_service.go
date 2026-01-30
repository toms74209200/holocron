package book

import (
	"context"

	"holocron/internal/book/domain"
)

type ListBooksInput struct {
	Q      *string
	Limit  *int
	Offset *int
}

type ListBooksOutput struct {
	Items []domain.BookItem
	Total int64
}

func ListBooks(
	ctx context.Context,
	queries *Queries,
	input ListBooksInput,
) (*ListBooksOutput, error) {
	keyword := domain.ToSearchKeyword(input.Q)
	pagination := domain.ToPagination(input.Limit, input.Offset)

	sources := []domain.BookListSource{
		SearchBooksSource(queries),
		ListAllBooksSource(queries),
	}
	items, total, err := domain.GetBookList(ctx, sources, keyword, pagination)
	if err != nil {
		return nil, err
	}

	return &ListBooksOutput{Items: items, Total: total}, nil
}
