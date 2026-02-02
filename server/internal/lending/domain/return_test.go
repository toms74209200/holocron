//go:build small

package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// When ValidateReturn with nil currentLending then returns ErrBookNotBorrowed
func TestValidateReturn_WithNilCurrentLending_ReturnsErrBookNotBorrowed(t *testing.T) {
	requestedBorrowerID := "borrower-1"
	actualBorrowerID := "borrower-1"

	err := ValidateReturn(nil, requestedBorrowerID, actualBorrowerID)

	if !errors.Is(err, ErrBookNotBorrowed) {
		t.Errorf("expected ErrBookNotBorrowed, got %v", err)
	}
}

// When ValidateReturn with different borrower IDs then returns error
func TestValidateReturn_WithDifferentBorrowerIDs_ReturnsError(t *testing.T) {
	currentLending := &CurrentLending{DueDate: time.Now()}
	requestedBorrowerID := "borrower-1"
	actualBorrowerID := "borrower-2"

	err := ValidateReturn(currentLending, requestedBorrowerID, actualBorrowerID)

	if !errors.Is(err, ErrNotBorrower) {
		t.Errorf("expected ErrNotBorrower, got %v", err)
	}
}

// When ValidateReturn with matching borrower IDs then returns nil
func TestValidateReturn_WithMatchingBorrowerIDs_ReturnsNil(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns nil for any matching borrower ID", prop.ForAll(
		func(borrowerID string, dueDate time.Time) bool {
			currentLending := &CurrentLending{DueDate: dueDate}
			err := ValidateReturn(currentLending, borrowerID, borrowerID)
			return err == nil
		},
		gen.AnyString(),
		gen.Time(),
	))
	properties.TestingRun(t)
}

// When ValidateReturn with different borrower IDs then always returns error
func TestValidateReturn_WithDifferentBorrowerIDs_AlwaysReturnsError(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns error for any different borrower IDs", prop.ForAll(
		func(requestedID string, actualID string, dueDate time.Time) bool {
			if requestedID == actualID {
				return true // Skip when IDs are equal
			}
			currentLending := &CurrentLending{DueDate: dueDate}
			err := ValidateReturn(currentLending, requestedID, actualID)
			return errors.Is(err, ErrNotBorrower)
		},
		gen.AnyString(),
		gen.AnyString(),
		gen.Time(),
	))
	properties.TestingRun(t)
}
