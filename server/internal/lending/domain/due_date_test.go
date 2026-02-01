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

// When CalculateDueDate with nil currentLending then returns due date from now
func TestCalculateDueDate_WithNilCurrentLending_ReturnsDueDateFromNow(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	dueDate, dueDays, err := CalculateDueDate(now, nil, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dueDays != 7 {
		t.Errorf("expected dueDays 7, got %d", dueDays)
	}
	expectedDueDate := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)
	if !dueDate.Equal(expectedDueDate) {
		t.Errorf("expected dueDate %v, got %v", expectedDueDate, dueDate)
	}
}

// When CalculateDueDate with currentLending then returns due date from current due date
func TestCalculateDueDate_WithCurrentLending_ReturnsDueDateFromCurrentDueDate(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	currentDueDate := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)
	currentLending := &CurrentLending{DueDate: currentDueDate}

	dueDate, dueDays, err := CalculateDueDate(now, nil, currentLending)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dueDays != 7 {
		t.Errorf("expected dueDays 7, got %d", dueDays)
	}
	expectedDueDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	if !dueDate.Equal(expectedDueDate) {
		t.Errorf("expected dueDate %v, got %v", expectedDueDate, dueDate)
	}
}

// When CalculateDueDate with custom dueDays then returns correct due date
func TestCalculateDueDate_WithCustomDueDays_ReturnsCorrectDueDate(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns correct due date", prop.ForAll(
		func(now time.Time, dueDaysValue int) bool {
			expectedDueDate := now.AddDate(0, 0, dueDaysValue)

			dueDate, dueDays, err := CalculateDueDate(now, &dueDaysValue, nil)

			return err == nil &&
				dueDays == dueDaysValue &&
				dueDate.Equal(expectedDueDate)
		},
		gen.Time(),
		gen.IntRange(1, 10000),
	))
	properties.TestingRun(t)
}

// When CalculateDueDate with zero dueDays then returns error
func TestCalculateDueDate_WithZeroDueDays_ReturnsError(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	dueDays := 0

	_, _, err := CalculateDueDate(now, &dueDays, nil)

	if !errors.Is(err, ErrInvalidDueDays) {
		t.Errorf("expected ErrInvalidDueDays, got %v", err)
	}
}

// When CalculateDueDate with negative dueDays then returns error
func TestCalculateDueDate_WithNegativeDueDays_ReturnsError(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns ErrInvalidDueDays", prop.ForAll(
		func(dueDays int) bool {
			now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
			_, _, err := CalculateDueDate(now, &dueDays, nil)
			return errors.Is(err, ErrInvalidDueDays)
		},
		gen.IntRange(-10000, -1),
	))
	properties.TestingRun(t)
}
