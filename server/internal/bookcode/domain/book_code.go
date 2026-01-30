package domain

import "errors"

var ErrInvalidCode = errors.New("code must not be empty")

type BookCode string

func ParseBookCode(s string) (BookCode, error) {
	if s == "" {
		return "", ErrInvalidCode
	}
	return BookCode(s), nil
}
