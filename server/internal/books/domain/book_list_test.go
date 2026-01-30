//go:build small

package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

var errNotMyResponsibility = errors.New("not my responsibility")

func TestGetBookList_WithFirstSourceSucceeds_ReturnsFirstSourceResult(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns first source result", prop.ForAll(
		func(total int64) bool {
			items := []BookItem{{ID: "first"}}
			sources := []BookListSource{
				func(ctx context.Context, keyword *SearchKeyword, pagination Pagination) ([]BookItem, int64, error) {
					return items, total, nil
				},
				func(ctx context.Context, keyword *SearchKeyword, pagination Pagination) ([]BookItem, int64, error) {
					return []BookItem{{ID: "second"}}, 999, nil
				},
			}

			result, resultTotal, err := GetBookList(context.Background(), sources, nil, Pagination{})
			return err == nil && len(result) == 1 && result[0].ID == "first" && resultTotal == total
		},
		gen.Int64(),
	))
	properties.TestingRun(t)
}

func TestGetBookList_WithFirstSourceFails_ReturnsSecondSourceResult(t *testing.T) {
	items := []BookItem{{ID: "second"}}
	sources := []BookListSource{
		func(ctx context.Context, keyword *SearchKeyword, pagination Pagination) ([]BookItem, int64, error) {
			return nil, 0, errNotMyResponsibility
		},
		func(ctx context.Context, keyword *SearchKeyword, pagination Pagination) ([]BookItem, int64, error) {
			return items, 1, nil
		},
	}

	result, total, err := GetBookList(context.Background(), sources, nil, Pagination{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].ID != "second" {
		t.Errorf("expected second source result, got %v", result)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
}

func TestGetBookList_WithAllSourcesFail_ReturnsError(t *testing.T) {
	sources := []BookListSource{
		func(ctx context.Context, keyword *SearchKeyword, pagination Pagination) ([]BookItem, int64, error) {
			return nil, 0, errNotMyResponsibility
		},
		func(ctx context.Context, keyword *SearchKeyword, pagination Pagination) ([]BookItem, int64, error) {
			return nil, 0, errNotMyResponsibility
		},
	}

	_, _, err := GetBookList(context.Background(), sources, nil, Pagination{})
	if err == nil {
		t.Error("expected error, got nil")
	}
}
