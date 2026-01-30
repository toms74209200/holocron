//go:build small

package domain

import (
	"errors"
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// When ParseBookCode with non-empty string then returns BookCode
func TestParseBookCode_WithNonEmptyString_ReturnsBookCode(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns BookCode with same value", prop.ForAll(
		func(s string) bool {
			code, err := ParseBookCode(s)
			return err == nil && string(code) == s
		},
		genNonEmptyString(),
	))
	properties.TestingRun(t)
}

// When ParseBookCode with empty string then returns error
func TestParseBookCode_WithEmptyString_ReturnsInvalidCodeError(t *testing.T) {
	_, err := ParseBookCode("")
	if !errors.Is(err, ErrInvalidCode) {
		t.Errorf("expected ErrInvalidCode, got %v", err)
	}
}

func genNonEmptyString() gopter.Gen {
	return gen.IntRange(1, 20).FlatMap(func(v interface{}) gopter.Gen {
		length := v.(int)
		return gen.SliceOfN(length, gen.NumChar()).Map(func(chars []rune) string {
			return string(chars)
		})
	}, reflect.TypeOf(""))
}
