//go:build small

package domain

import "testing"

func TestParseBookAuthors_WithValidAuthors_ReturnsBookAuthors(t *testing.T) {
	input := []string{"Author1", "Author2"}

	authors, err := ParseBookAuthors(input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(authors) != 2 {
		t.Errorf("expected 2 authors, got %d", len(authors))
	}
	if authors[0] != "Author1" {
		t.Errorf("expected Author1, got %s", authors[0])
	}
}

func TestParseBookAuthors_WithEmptySlice_ReturnsError(t *testing.T) {
	_, err := ParseBookAuthors([]string{})

	if err != ErrInvalidAuthors {
		t.Errorf("expected ErrInvalidAuthors, got %v", err)
	}
}

func TestParseBookAuthors_WithNilSlice_ReturnsError(t *testing.T) {
	_, err := ParseBookAuthors(nil)

	if err != ErrInvalidAuthors {
		t.Errorf("expected ErrInvalidAuthors, got %v", err)
	}
}
