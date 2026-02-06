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

// When ParseDeleteReason with valid reason then returns DeleteReason
func TestParseDeleteReason_WithValidReason_ReturnsDeleteReason(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns DeleteReason with same value", prop.ForAll(
		func(s string) bool {
			reason, err := ParseDeleteReason(s)
			return err == nil && string(reason) == s
		},
		genValidDeleteReason(),
	))
	properties.TestingRun(t)
}

// When ParseDeleteReason with invalid reason then returns error
func TestParseDeleteReason_WithInvalidReason_ReturnsError(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns ErrInvalidDeleteReason", prop.ForAll(
		func(s string) bool {
			_, err := ParseDeleteReason(s)
			return errors.Is(err, ErrInvalidDeleteReason)
		},
		genInvalidDeleteReason(),
	))
	properties.TestingRun(t)
}

// When ParseDeleteReason with empty string then returns error
func TestParseDeleteReason_WithEmptyString_ReturnsError(t *testing.T) {
	_, err := ParseDeleteReason("")
	if !errors.Is(err, ErrInvalidDeleteReason) {
		t.Errorf("expected ErrInvalidDeleteReason, got %v", err)
	}
}

func genValidDeleteReason() gopter.Gen {
	validReasons := []string{
		string(DeleteReasonTransfer),
		string(DeleteReasonDisposal),
		string(DeleteReasonLost),
		string(DeleteReasonOther),
	}
	return gen.OneConstOf(
		DeleteReasonTransfer,
		DeleteReasonDisposal,
		DeleteReasonLost,
		DeleteReasonOther,
	).Map(func(v DeleteReason) string {
		return string(v)
	}).SuchThat(func(v string) bool {
		for _, valid := range validReasons {
			if v == valid {
				return true
			}
		}
		return false
	})
}

func genInvalidDeleteReason() gopter.Gen {
	return gen.OneConstOf(
		"invalid",
		"unknown",
		"delete",
		"remove",
		"Transfer",
		"TRANSFER",
		"disposal ",
		" disposal",
	).FlatMap(func(v interface{}) gopter.Gen {
		return gen.Const(v.(string))
	}, reflect.TypeOf(""))
}
