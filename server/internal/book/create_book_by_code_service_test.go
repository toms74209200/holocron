//go:build medium

package book

import (
	"context"
	"testing"

	"holocron/internal/book/domain"
)

func TestCreateBookByCode_WithValidCode_ReturnsBook(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	googleFetcher, err := NewGoogleBooksFetcher()
	if err != nil {
		t.Fatal(err)
	}
	openBDFetcher, err := NewOpenBDFetcher()
	if err != nil {
		t.Fatal(err)
	}
	sources := []domain.BookInfoSource{
		DBCacheSource(queries),
		ExternalAPISource(googleFetcher.Fetch, domain.BookInfoFromGoogleBooks),
		ExternalAPISource(openBDFetcher.Fetch, domain.BookInfoFromOpenBD),
	}

	output, err := CreateBookByCode(context.Background(), queries, sources, CreateBookByCodeInput{
		Code: "9784873115658",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Title != "リーダブルコード" {
		t.Errorf("expected title %q, got %q", "リーダブルコード", output.Title)
	}
}

func TestCreateBookByCode_WithEmptyCode_ReturnsInvalidCodeError(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)

	_, err := CreateBookByCode(context.Background(), queries, nil, CreateBookByCodeInput{
		Code: "",
	})

	if err != ErrInvalidCode {
		t.Errorf("expected ErrInvalidCode, got %v", err)
	}
}

func TestCreateBookByCode_WithSameCode_ReturnsBookWithSameInfo(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	googleFetcher, err := NewGoogleBooksFetcher()
	if err != nil {
		t.Fatal(err)
	}
	openBDFetcher, err := NewOpenBDFetcher()
	if err != nil {
		t.Fatal(err)
	}
	sources := []domain.BookInfoSource{
		DBCacheSource(queries),
		ExternalAPISource(googleFetcher.Fetch, domain.BookInfoFromGoogleBooks),
		ExternalAPISource(openBDFetcher.Fetch, domain.BookInfoFromOpenBD),
	}
	first, _ := CreateBookByCode(context.Background(), queries, sources, CreateBookByCodeInput{
		Code: "9784873115658",
	})

	second, err := CreateBookByCode(context.Background(), queries, sources, CreateBookByCodeInput{
		Code: "9784873115658",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if first.ID == second.ID {
		t.Error("expected different book IDs")
	}
	if second.Title != first.Title {
		t.Errorf("expected same title %q, got %q", first.Title, second.Title)
	}
}

