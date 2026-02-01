package domain

import (
	"errors"
	"time"
)

const DefaultDueDays = 7

var ErrInvalidDueDays = errors.New("due days must be at least 1")

func CalculateDueDate(now time.Time, requestedDueDays *int, currentLending *CurrentLending) (dueDate time.Time, dueDays int, err error) {
	dueDays = DefaultDueDays
	if requestedDueDays != nil {
		if *requestedDueDays < 1 {
			return time.Time{}, 0, ErrInvalidDueDays
		}
		dueDays = *requestedDueDays
	}

	baseDate := now
	if currentLending != nil {
		baseDate = currentLending.DueDate
	}

	dueDate = baseDate.AddDate(0, 0, dueDays)
	return dueDate, dueDays, nil
}
