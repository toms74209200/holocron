//go:build small

package domain

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestBookInfoFromGoogleBooks_WithValidResponse_ReturnsSameValues(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns same title", prop.ForAll(
		func(title string, author string, publisher string) bool {
			body := genGoogleBooksJSON(title, author, publisher)
			info, err := BookInfoFromGoogleBooks(body)
			return err == nil && info.Title == title && info.Publisher == publisher
		},
		gen.AnyString(),
		gen.AnyString(),
		gen.AnyString(),
	))
	properties.TestingRun(t)
}

func TestBookInfoFromGoogleBooks_WithEmptyItems_ReturnsNotFoundError(t *testing.T) {
	body := []byte(`{"totalItems": 0}`)

	_, err := BookInfoFromGoogleBooks(body)

	if err != ErrBookNotFound {
		t.Errorf("expected ErrBookNotFound, got %v", err)
	}
}

func TestBookInfoFromGoogleBooks_WithInvalidJSON_ReturnsError(t *testing.T) {
	body := []byte(`invalid json`)

	_, err := BookInfoFromGoogleBooks(body)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestBookInfoFromOpenBD_WithValidResponse_ReturnsSameValues(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns same title and normalized date", prop.ForAll(
		func(title string, publisher string, year, month, day int) bool {
			if title == "" {
				return true
			}
			pubdate := fmt.Sprintf("%04d%02d%02d", year, month, day)
			expected := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
			body := genOpenBDJSON(title, publisher, pubdate)
			info, err := BookInfoFromOpenBD(body)
			return err == nil && info.Title == title && info.Publisher == publisher && info.PublishedDate == expected
		},
		gen.AnyString(),
		gen.AnyString(),
		gen.IntRange(1900, 2100),
		gen.IntRange(1, 12),
		gen.IntRange(1, 28),
	))
	properties.TestingRun(t)
}

func TestBookInfoFromOpenBD_WithNullResult_ReturnsNotFoundError(t *testing.T) {
	body := []byte(`[null]`)

	_, err := BookInfoFromOpenBD(body)

	if err != ErrBookNotFound {
		t.Errorf("expected ErrBookNotFound, got %v", err)
	}
}

func TestBookInfoFromOpenBD_WithEmptyTitle_ReturnsNotFoundError(t *testing.T) {
	body := []byte(`[{"summary": {"title": ""}}]`)

	_, err := BookInfoFromOpenBD(body)

	if err != ErrBookNotFound {
		t.Errorf("expected ErrBookNotFound, got %v", err)
	}
}

func genGoogleBooksJSON(title, author, publisher string) []byte {
	resp := map[string]any{
		"totalItems": 1,
		"items": []map[string]any{
			{
				"volumeInfo": map[string]any{
					"title":     title,
					"authors":   []string{author},
					"publisher": publisher,
				},
			},
		},
	}
	b, _ := json.Marshal(resp)
	return b
}

func genOpenBDJSON(title, publisher, pubdate string) []byte {
	resp := []map[string]any{
		{
			"summary": map[string]any{
				"title":     title,
				"publisher": publisher,
				"pubdate":   pubdate,
			},
		},
	}
	b, _ := json.Marshal(resp)
	return b
}

func init() {
	_ = reflect.TypeOf("")
}
