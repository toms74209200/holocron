package domain

import "errors"

var ErrInvalidName = errors.New("name must be 1-50 characters")

type UserName string

func ParseUserName(s string) (UserName, error) {
	if len(s) < 1 || len(s) > 50 {
		return "", ErrInvalidName
	}
	return UserName(s), nil
}
