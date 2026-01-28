package user

import (
	"context"
	"errors"
	"time"

	"holocron/internal/user/domain"

	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUserName   = errors.New("invalid user name")
	ErrTokenCreation     = errors.New("failed to create token")
)

type FirebaseAuth interface {
	CreateCustomToken(ctx context.Context, uid string) (string, error)
}

type CreateUserInput struct {
	Name *string
}

type CreateUserOutput struct {
	ID          string
	Name        string
	CustomToken string
	CreatedAt   time.Time
}

func CreateUser(ctx context.Context, queries *Queries, firebaseAuth FirebaseAuth, input CreateUserInput) (*CreateUserOutput, error) {
	userID := uuid.New().String()

	var userName domain.UserName
	var err error
	if input.Name == nil || *input.Name == "" {
		userName = domain.UserName("ユーザー" + userID[:4])
	} else {
		userName, err = domain.ParseUserName(*input.Name)
		if err != nil {
			return nil, ErrInvalidUserName
		}
	}

	cnt, err := queries.CountUserByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, ErrUserAlreadyExists
	}

	customToken, err := firebaseAuth.CreateCustomToken(ctx, userID)
	if err != nil {
		return nil, ErrTokenCreation
	}

	now := time.Now().UTC()
	err = queries.InsertUserEvent(ctx, InsertUserEventParams{
		EventID:    uuid.New().String(),
		UserID:     userID,
		EventType:  "created",
		Name:       string(userName),
		OccurredAt: now.Format(time.RFC3339),
	})
	if err != nil {
		return nil, err
	}

	return &CreateUserOutput{
		ID:          userID,
		Name:        string(userName),
		CustomToken: customToken,
		CreatedAt:   now,
	}, nil
}
