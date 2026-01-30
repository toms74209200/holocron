//go:build small

package domain

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestToPagination_WithValidLimit_ReturnsSameLimit(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns same limit when valid", prop.ForAll(
		func(limit int) bool {
			p := ToPagination(&limit, nil)
			return p.Limit() == limit
		},
		gen.IntRange(1, 100),
	))
	properties.TestingRun(t)
}

func TestToPagination_WithInvalidLimit_ReturnsDefault(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns default 20 when limit <= 0", prop.ForAll(
		func(limit int) bool {
			p := ToPagination(&limit, nil)
			return p.Limit() == 20
		},
		gen.IntRange(-100, 0),
	))
	properties.Property("returns default 20 when limit > 100", prop.ForAll(
		func(limit int) bool {
			p := ToPagination(&limit, nil)
			return p.Limit() == 20
		},
		gen.IntRange(101, 1000),
	))
	properties.TestingRun(t)
}

func TestToPagination_WithValidOffset_ReturnsSameOffset(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns same offset when valid", prop.ForAll(
		func(offset int) bool {
			p := ToPagination(nil, &offset)
			return p.Offset() == offset
		},
		gen.IntRange(0, 1000),
	))
	properties.TestingRun(t)
}

func TestToPagination_WithNegativeOffset_ReturnsZero(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("returns 0 when offset < 0", prop.ForAll(
		func(offset int) bool {
			p := ToPagination(nil, &offset)
			return p.Offset() == 0
		},
		gen.IntRange(-1000, -1),
	))
	properties.TestingRun(t)
}

func TestToPagination_WithNil_ReturnsDefaults(t *testing.T) {
	p := ToPagination(nil, nil)
	if p.Limit() != 20 {
		t.Errorf("expected limit 20, got %d", p.Limit())
	}
	if p.Offset() != 0 {
		t.Errorf("expected offset 0, got %d", p.Offset())
	}
}
