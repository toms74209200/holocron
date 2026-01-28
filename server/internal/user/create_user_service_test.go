//go:build medium

package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	_, err = db.Exec(`
		CREATE TABLE user_events (
			event_id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			event_type TEXT NOT NULL,
			name TEXT NOT NULL,
			occurred_at TEXT NOT NULL
		);
		CREATE INDEX idx_user_events_user_id ON user_events(user_id);
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

type fakeFirebaseAuth struct {
	token string
	err   error
}

func (f *fakeFirebaseAuth) CreateCustomToken(_ context.Context, _ string) (string, error) {
	return f.token, f.err
}

// When Execute with new user then returns CreateUserOutput
func TestCreateUserService_Execute_WithNewUser_ReturnsOutput(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	firebaseAuth := &fakeFirebaseAuth{token: "test-token", err: nil}
	service := NewCreateUserService(queries, firebaseAuth)
	ctx := context.Background()

	name := "TestUser"
	input := CreateUserInput{Name: &name}

	output, err := service.Execute(ctx, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.ID == "" {
		t.Error("expected ID to be set")
	}
	if output.Name != name {
		t.Errorf("expected Name %s, got %s", name, output.Name)
	}
	if output.CustomToken != "test-token" {
		t.Errorf("expected CustomToken test-token, got %s", output.CustomToken)
	}
	if output.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

// When Execute with nil name then generates name automatically
func TestCreateUserService_Execute_WithNilName_GeneratesName(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	firebaseAuth := &fakeFirebaseAuth{token: "test-token", err: nil}
	service := NewCreateUserService(queries, firebaseAuth)
	ctx := context.Background()

	input := CreateUserInput{Name: nil}

	output, err := service.Execute(ctx, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Name == "" {
		t.Error("expected Name to be generated")
	}
	if len(output.Name) < 5 {
		t.Errorf("expected generated name to have prefix, got %s", output.Name)
	}
}

// When Execute with invalid name then returns ErrInvalidUserName
func TestCreateUserService_Execute_WithInvalidName_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	firebaseAuth := &fakeFirebaseAuth{token: "test-token", err: nil}
	service := NewCreateUserService(queries, firebaseAuth)
	ctx := context.Background()

	longName := "this name is way too long and exceeds the fifty character limit for user names"
	input := CreateUserInput{Name: &longName}

	_, err := service.Execute(ctx, input)

	if !errors.Is(err, ErrInvalidUserName) {
		t.Errorf("expected ErrInvalidUserName, got %v", err)
	}
}

// When Execute with firebase error then returns ErrTokenCreation
func TestCreateUserService_Execute_WithFirebaseError_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	queries := New(db)
	firebaseAuth := &fakeFirebaseAuth{token: "", err: errors.New("firebase error")}
	service := NewCreateUserService(queries, firebaseAuth)
	ctx := context.Background()

	name := "TestUser"
	input := CreateUserInput{Name: &name}

	_, err := service.Execute(ctx, input)

	if !errors.Is(err, ErrTokenCreation) {
		t.Errorf("expected ErrTokenCreation, got %v", err)
	}
}
