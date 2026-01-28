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

type CreateUserService struct {
	queries      *Queries
	firebaseAuth FirebaseAuth
}

func NewCreateUserService(queries *Queries, firebaseAuth FirebaseAuth) *CreateUserService {
	return &CreateUserService{
		queries:      queries,
		firebaseAuth: firebaseAuth,
	}
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

func (s *CreateUserService) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
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

	cnt, err := s.queries.CountUserByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, ErrUserAlreadyExists
	}

	now := time.Now().UTC()
	err = s.queries.InsertUserEvent(ctx, InsertUserEventParams{
		EventID:    uuid.New().String(),
		UserID:     userID,
		EventType:  "created",
		Name:       string(userName),
		OccurredAt: now.Format(time.RFC3339),
	})
	if err != nil {
		return nil, err
	}

	customToken, err := s.firebaseAuth.CreateCustomToken(ctx, userID)
	if err != nil {
		return nil, ErrTokenCreation
	}

	return &CreateUserOutput{
		ID:          userID,
		Name:        string(userName),
		CustomToken: customToken,
		CreatedAt:   now,
	}, nil
}
