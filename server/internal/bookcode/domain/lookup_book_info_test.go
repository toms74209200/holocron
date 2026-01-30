package domain

import (
	"context"
	"errors"
	"testing"

	book "holocron/internal/book/domain"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestLookupBookInfo_WithSuccessAtIndex_ReturnsSuccessValue(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("returns value from first successful source", prop.ForAll(
		func(n int, successIndex int, title string) bool {
			if n == 0 || title == "" {
				return true
			}
			successIndex = successIndex % n

			sources := make([]BookInfoSource, n)
			for i := range sources {
				if i < successIndex {
					sources[i] = func(ctx context.Context, code string) (*book.BookInfo, error) {
						return nil, errors.New("fail")
					}
				} else if i == successIndex {
					sources[i] = func(ctx context.Context, code string) (*book.BookInfo, error) {
						return &book.BookInfo{Title: title}, nil
					}
				} else {
					sources[i] = func(ctx context.Context, code string) (*book.BookInfo, error) {
						return &book.BookInfo{Title: "wrong"}, nil
					}
				}
			}

			info, err := LookupBookInfo(context.Background(), sources, "code")
			return err == nil && info.Title == title
		},
		gen.IntRange(1, 10),
		gen.IntRange(0, 9),
		gen.AnyString().SuchThat(func(s string) bool { return s != "" && s != "wrong" }),
	))

	properties.TestingRun(t)
}

func TestLookupBookInfo_WithAllFailing_ReturnsNotFoundError(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("returns ErrBookNotFound when all sources fail", prop.ForAll(
		func(n int) bool {
			sources := make([]BookInfoSource, n)
			for i := range sources {
				sources[i] = func(ctx context.Context, code string) (*book.BookInfo, error) {
					return nil, errors.New("fail")
				}
			}

			_, err := LookupBookInfo(context.Background(), sources, "code")
			return err == book.ErrBookNotFound
		},
		gen.IntRange(0, 10),
	))

	properties.TestingRun(t)
}
