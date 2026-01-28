package domain

import "errors"

var ErrInvalidTitle = errors.New("title must be 1-200 characters")

type BookTitle string

func ParseBookTitle(s string) (BookTitle, error) {
	if len(s) < 1 || len(s) > 200 {
		return "", ErrInvalidTitle
	}
	return BookTitle(s), nil
}
