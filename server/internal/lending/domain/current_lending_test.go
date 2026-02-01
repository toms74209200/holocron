//go:build small

package domain

import (
	"testing"
	"time"
)

// When ParseCurrentLending with latestDueDate then returns CurrentLending with that due date
func TestParseCurrentLending_WithLatestDueDate_ReturnsCurrentLendingWithThatDueDate(t *testing.T) {
	latestDueDate := "2024-01-15T00:00:00Z"
	borrowedDueDate := "2024-01-08T00:00:00Z"
	borrowedAt := "2024-01-01T00:00:00Z"

	currentLending, err := ParseCurrentLending(&latestDueDate, &borrowedDueDate, borrowedAt)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedDueDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	if !currentLending.DueDate.Equal(expectedDueDate) {
		t.Errorf("expected DueDate %v, got %v", expectedDueDate, currentLending.DueDate)
	}
}

// When ParseCurrentLending with nil latestDueDate then returns CurrentLending with borrowedDueDate
func TestParseCurrentLending_WithNilLatestDueDate_ReturnsCurrentLendingWithBorrowedDueDate(t *testing.T) {
	borrowedDueDate := "2024-01-08T00:00:00Z"
	borrowedAt := "2024-01-01T00:00:00Z"

	currentLending, err := ParseCurrentLending(nil, &borrowedDueDate, borrowedAt)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedDueDate := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)
	if !currentLending.DueDate.Equal(expectedDueDate) {
		t.Errorf("expected DueDate %v, got %v", expectedDueDate, currentLending.DueDate)
	}
}

// When ParseCurrentLending with nil latestDueDate and nil borrowedDueDate then calculates from borrowedAt
func TestParseCurrentLending_WithNilLatestDueDateAndNilBorrowedDueDate_CalculatesFromBorrowedAt(t *testing.T) {
	borrowedAt := "2024-01-01T00:00:00Z"

	currentLending, err := ParseCurrentLending(nil, nil, borrowedAt)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedDueDate := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)
	if !currentLending.DueDate.Equal(expectedDueDate) {
		t.Errorf("expected DueDate %v, got %v", expectedDueDate, currentLending.DueDate)
	}
}

// When ParseCurrentLending with invalid latestDueDate then returns error
func TestParseCurrentLending_WithInvalidLatestDueDate_ReturnsError(t *testing.T) {
	latestDueDate := "invalid-date"
	borrowedDueDate := "2024-01-08T00:00:00Z"
	borrowedAt := "2024-01-01T00:00:00Z"

	_, err := ParseCurrentLending(&latestDueDate, &borrowedDueDate, borrowedAt)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

// When ParseCurrentLending with invalid borrowedDueDate then returns error
func TestParseCurrentLending_WithInvalidBorrowedDueDate_ReturnsError(t *testing.T) {
	borrowedDueDate := "invalid-date"
	borrowedAt := "2024-01-01T00:00:00Z"

	_, err := ParseCurrentLending(nil, &borrowedDueDate, borrowedAt)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

// When ParseCurrentLending with invalid borrowedAt then returns error
func TestParseCurrentLending_WithInvalidBorrowedAt_ReturnsError(t *testing.T) {
	borrowedAt := "invalid-date"

	_, err := ParseCurrentLending(nil, nil, borrowedAt)

	if err == nil {
		t.Error("expected error, got nil")
	}
}
