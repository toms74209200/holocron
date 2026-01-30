//go:build small

package domain

import (
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestToSearchKeyword_WithNonEmptyString_ReturnsKeyword(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns keyword with same value", prop.ForAll(
		func(s string) bool {
			k := ToSearchKeyword(&s)
			return k != nil && string(*k) == s
		},
		genNonEmptyString(),
	))
	properties.TestingRun(t)
}

func TestToSearchKeyword_WithNil_ReturnsNil(t *testing.T) {
	k := ToSearchKeyword(nil)
	if k != nil {
		t.Errorf("expected nil, got %v", k)
	}
}

func TestToSearchKeyword_WithEmptyString_ReturnsNil(t *testing.T) {
	q := ""
	k := ToSearchKeyword(&q)
	if k != nil {
		t.Errorf("expected nil, got %v", k)
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
