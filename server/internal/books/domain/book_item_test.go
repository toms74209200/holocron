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
			expectedID := bookID
			expectedTitle := title
			expectedAuthors := authors
			expectedStatus := "available"
			expectedCreatedAt := occurredAt.Truncate(time.Second)
			authorsJSON, _ := json.Marshal(authors)

			actual, err := BookItemFromRow(
				bookID,
				nil,
				title,
				string(authorsJSON),
				nil,
				nil,
				nil,
				occurredAt.Format(time.RFC3339),
				nil,
				nil,
				nil,
			)

			if err != nil {
				return false
			}
			return actual.ID == expectedID &&
				actual.Title == expectedTitle &&
				len(actual.Authors) == len(expectedAuthors) &&
				actual.Status == expectedStatus &&
				actual.Borrower == nil &&
				actual.CreatedAt.Equal(expectedCreatedAt)
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
	properties := gopter.NewProperties(nil)
	properties.Property("returns error when authorsJSON is invalid", prop.ForAll(
		func(bookID, title string, occurredAt time.Time) bool {
			invalidAuthorsJSON := "invalid-json"

			_, err := BookItemFromRow(
				bookID,
				nil,
				title,
				invalidAuthorsJSON,
				nil,
				nil,
				nil,
				occurredAt.Format(time.RFC3339),
				nil,
				nil,
				nil,
			)

			return err == ErrInvalidBookRow
		},
		gen.AnyString(),
		gen.AnyString(),
		gen.Time(),
	))
	properties.TestingRun(t)
}

// When BookItemFromRow with invalid occurredAt then returns error
func TestBookItemFromRow_WithInvalidOccurredAt_ReturnsError(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns error when occurredAt is invalid", prop.ForAll(
		func(bookID, title string, authors []string) bool {
			authorsJSON, _ := json.Marshal(authors)
			invalidOccurredAt := "invalid-date"

			_, err := BookItemFromRow(
				bookID,
				nil,
				title,
				string(authorsJSON),
				nil,
				nil,
				nil,
				invalidOccurredAt,
				nil,
				nil,
				nil,
			)

			return err == ErrInvalidBookRow
		},
		gen.AnyString(),
		gen.AnyString(),
		gen.SliceOf(gen.AnyString()),
	))
	properties.TestingRun(t)
}

// When BookItemFromRow with borrower info then returns borrowed status
func TestBookItemFromRow_WithBorrowerInfo_ReturnsBorrowedStatus(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns borrowed status with borrower info", prop.ForAll(
		func(bookID, title, borrowerID, borrowerName string, authors []string, occurredAt, borrowedAt time.Time) bool {
			expectedID := bookID
			expectedTitle := title
			expectedStatus := "borrowed"
			expectedBorrowerID := borrowerID
			expectedBorrowerName := borrowerName
			expectedBorrowedAt := borrowedAt.Truncate(time.Second)
			authorsJSON, _ := json.Marshal(authors)
			borrowedAtStr := borrowedAt.Format(time.RFC3339)

			actual, err := BookItemFromRow(
				bookID,
				nil,
				title,
				string(authorsJSON),
				nil,
				nil,
				nil,
				occurredAt.Format(time.RFC3339),
				&borrowerID,
				&borrowerName,
				&borrowedAtStr,
			)

			if err != nil {
				return false
			}
			return actual.ID == expectedID &&
				actual.Title == expectedTitle &&
				actual.Status == expectedStatus &&
				actual.Borrower != nil &&
				actual.Borrower.ID == expectedBorrowerID &&
				actual.Borrower.Name == expectedBorrowerName &&
				actual.Borrower.BorrowedAt.Equal(expectedBorrowedAt)
		},
		gen.AnyString(),
		gen.AnyString(),
		gen.AnyString(),
		gen.AnyString(),
		gen.SliceOf(gen.AnyString()),
		gen.Time(),
		gen.Time(),
	))
	properties.TestingRun(t)
}

// When BookItemFromRow with partial borrower info then returns available status
func TestBookItemFromRow_WithPartialBorrowerInfo_ReturnsAvailableStatus(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns available status when borrower info is incomplete", prop.ForAll(
		func(bookID, title, borrowerID string, authors []string, occurredAt time.Time) bool {
			expectedStatus := "available"
			authorsJSON, _ := json.Marshal(authors)

			actual, err := BookItemFromRow(
				bookID,
				nil,
				title,
				string(authorsJSON),
				nil,
				nil,
				nil,
				occurredAt.Format(time.RFC3339),
				&borrowerID,
				nil,
				nil,
			)

			if err != nil {
				return false
			}
			return actual.Status == expectedStatus && actual.Borrower == nil
		},
		gen.AnyString(),
		gen.AnyString(),
		gen.AnyString(),
		gen.SliceOf(gen.AnyString()),
		gen.Time(),
	))
	properties.TestingRun(t)
}

// When BookItemFromRow with invalid borrowedAt then returns error
func TestBookItemFromRow_WithInvalidBorrowedAt_ReturnsError(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns error when borrowedAt is invalid", prop.ForAll(
		func(bookID, title, borrowerID, borrowerName string, authors []string, occurredAt time.Time) bool {
			authorsJSON, _ := json.Marshal(authors)
			invalidBorrowedAt := "invalid-date"

			_, err := BookItemFromRow(
				bookID,
				nil,
				title,
				string(authorsJSON),
				nil,
				nil,
				nil,
				occurredAt.Format(time.RFC3339),
				&borrowerID,
				&borrowerName,
				&invalidBorrowedAt,
			)

			return err == ErrInvalidBookRow
		},
		gen.AnyString(),
		gen.AnyString(),
		gen.AnyString(),
		gen.AnyString(),
		gen.SliceOf(gen.AnyString()),
		gen.Time(),
	))
	properties.TestingRun(t)
}
