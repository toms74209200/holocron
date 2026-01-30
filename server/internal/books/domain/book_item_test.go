//go:build small

package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// When BookItemFromRow with valid row then returns BookItem with same values
func TestBookItemFromRow_WithValidRow_ReturnsBookItem(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns BookItem with same values", prop.ForAll(
		func(bookID, title string, authors []string, occurredAt time.Time) bool {
			authorsJSON, _ := json.Marshal(authors)
			item, err := BookItemFromRow(
				bookID,
				nil,
				title,
				string(authorsJSON),
				nil,
				nil,
				nil,
				occurredAt.Format(time.RFC3339),
			)
			if err != nil {
				return false
			}
			return item.ID == bookID &&
				item.Title == title &&
				len(item.Authors) == len(authors) &&
				item.Status == "available"
		},
		gen.AnyString(),
		gen.AnyString(),
		gen.SliceOf(gen.AnyString()),
		gen.Time(),
	))
	properties.TestingRun(t)
}

// When BookItemFromRow with invalid authorsJSON then returns error
func TestBookItemFromRow_WithInvalidAuthorsJSON_ReturnsError(t *testing.T) {
	_, err := BookItemFromRow(
		"book-id",
		nil,
		"title",
		"invalid-json",
		nil,
		nil,
		nil,
		time.Now().Format(time.RFC3339),
	)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

// When BookItemFromRow with invalid occurredAt then returns error
func TestBookItemFromRow_WithInvalidOccurredAt_ReturnsError(t *testing.T) {
	authorsJSON, _ := json.Marshal([]string{"author"})
	_, err := BookItemFromRow(
		"book-id",
		nil,
		"title",
		string(authorsJSON),
		nil,
		nil,
		nil,
		"invalid-date",
	)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
