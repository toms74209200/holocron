//go:build small

package domain

import "testing"

func TestParseBookTitle_WithValidString_ReturnsBookTitle(t *testing.T) {
	title, err := ParseBookTitle("Test Book Title")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(title) != "Test Book Title" {
		t.Errorf("expected 'Test Book Title', got %s", title)
	}
}

func TestParseBookTitle_WithEmptyString_ReturnsError(t *testing.T) {
	_, err := ParseBookTitle("")

	if err != ErrInvalidTitle {
		t.Errorf("expected ErrInvalidTitle, got %v", err)
	}
}

func TestParseBookTitle_WithStringOver200Chars_ReturnsError(t *testing.T) {
	longTitle := make([]byte, 201)
	for i := range longTitle {
		longTitle[i] = 'a'
	}

	_, err := ParseBookTitle(string(longTitle))

	if err != ErrInvalidTitle {
		t.Errorf("expected ErrInvalidTitle, got %v", err)
	}
}
