package domain

import "errors"

var ErrInvalidAuthors = errors.New("authors must have at least one author")

type BookAuthors []string

func ParseBookAuthors(authors []string) (BookAuthors, error) {
	if len(authors) == 0 {
		return nil, ErrInvalidAuthors
	}
	return BookAuthors(authors), nil
}
