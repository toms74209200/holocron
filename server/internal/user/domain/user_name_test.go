//go:build small

package domain

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// When ParseUserName with valid string (1-50 chars) then returns UserName
func TestParseUserName_WithValidString_ReturnsUserName(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns UserName with same value", prop.ForAll(
		func(s string) bool {
			name, err := ParseUserName(s)
			return err == nil && string(name) == s
		},
		genValidUserName(),
	))
	properties.TestingRun(t)
}

// When ParseUserName with empty string then returns error
func TestParseUserName_WithEmptyString_ReturnsError(t *testing.T) {
	_, err := ParseUserName("")
	if !errors.Is(err, ErrInvalidName) {
		t.Errorf("expected ErrInvalidName, got %v", err)
	}
}

// When ParseUserName with string over 50 chars then returns error
func TestParseUserName_WithStringOver50Chars_ReturnsError(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns error", prop.ForAll(
		func(s string) bool {
			_, err := ParseUserName(s)
			return errors.Is(err, ErrInvalidName)
		},
		genInvalidLongUserName(),
	))
	properties.TestingRun(t)
}

func genValidUserName() gopter.Gen {
	return gen.IntRange(1, 50).FlatMap(func(v interface{}) gopter.Gen {
		length := v.(int)
		return gen.SliceOfN(length, gen.AlphaChar()).Map(func(chars []rune) string {
			return string(chars)
		})
	}, reflect.TypeOf(""))
}

func genInvalidLongUserName() gopter.Gen {
	return gen.IntRange(51, 100).FlatMap(func(v interface{}) gopter.Gen {
		length := v.(int)
		return gen.Const(strings.Repeat("a", length))
	}, reflect.TypeOf(""))
}
