package domain

import (
	"time"
)

type CurrentLending struct {
	DueDate time.Time
}

func ParseCurrentLending(latestDueDate, borrowedDueDate *string, borrowedAt string) (*CurrentLending, error) {
	var dueDate time.Time
	var err error

	if latestDueDate != nil {
		dueDate, err = time.Parse(time.RFC3339, *latestDueDate)
		if err != nil {
			return nil, err
		}
	} else if borrowedDueDate != nil {
		dueDate, err = time.Parse(time.RFC3339, *borrowedDueDate)
		if err != nil {
			return nil, err
		}
	} else {
		borrowedAtTime, err := time.Parse(time.RFC3339, borrowedAt)
		if err != nil {
			return nil, err
		}
		dueDate = borrowedAtTime.AddDate(0, 0, DefaultDueDays)
	}

	return &CurrentLending{DueDate: dueDate}, nil
}
