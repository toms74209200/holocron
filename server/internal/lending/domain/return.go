package domain

import (
	"errors"
)

var (
	ErrBookNotBorrowed = errors.New("book is not currently borrowed")
	ErrNotBorrower     = errors.New("only the borrower can return this book")
)

func ValidateReturn(currentLending *CurrentLending, requestedBorrowerID string, actualBorrowerID string) error {
	if currentLending == nil {
		return ErrBookNotBorrowed
	}

	if actualBorrowerID != requestedBorrowerID {
		return ErrNotBorrower
	}

	return nil
}
