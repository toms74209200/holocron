package auth

import (
	"context"
	"errors"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type FirebaseAuth interface {
	CreateCustomToken(ctx context.Context, uid string) (string, error)
	VerifyIDToken(ctx context.Context, idToken string) (string, error)
}

type firebaseAuth struct {
	client *auth.Client
}

func NewFirebaseAuth(ctx context.Context) (FirebaseAuth, error) {
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		projectID = "holocron"
	}

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID})
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &firebaseAuth{client: client}, nil
}

func (f *firebaseAuth) CreateCustomToken(ctx context.Context, uid string) (string, error) {
	return f.client.CustomToken(ctx, uid)
}

func (f *firebaseAuth) VerifyIDToken(ctx context.Context, idToken string) (string, error) {
	token, err := f.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", ErrInvalidToken
	}
	return token.UID, nil
}
