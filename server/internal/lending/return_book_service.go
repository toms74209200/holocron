package lending

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"holocron/internal/lending/domain"

	"github.com/google/uuid"
)

var (
	ErrBookNotBorrowed = errors.New("book is not currently borrowed")
	ErrNotBorrower     = errors.New("only the borrower can return this book")
)

type ReturnBookInput struct {
	BookID      string
	RequesterID string
}

type ReturnBookOutput struct {
	LendingID  string
	BookID     string
	BorrowerID string
	ReturnedAt time.Time
}

type ReturnBookService struct {
	lendingQueries *Queries
	bookQueries    BookQueries
	now            func() time.Time
}

func NewReturnBookService(lendingQueries *Queries, bookQueries BookQueries) *ReturnBookService {
	return &ReturnBookService{
		lendingQueries: lendingQueries,
		bookQueries:    bookQueries,
		now:            func() time.Time { return time.Now().UTC() },
	}
}

func (s *ReturnBookService) ReturnBook(ctx context.Context, input ReturnBookInput) (*ReturnBookOutput, error) {
	now := s.now()

	count, err := s.bookQueries.CountBookByBookId(ctx, input.BookID)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, ErrBookNotFound
	}

	currentLendingRow, err := s.lendingQueries.GetCurrentLending(ctx, input.BookID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBookNotBorrowed
		}
		return nil, err
	}

	latestDueDate, err := s.lendingQueries.GetLatestDueDate(ctx, currentLendingRow.LendingID)
	if err != nil {
		return nil, err
	}

	var latestDueDatePtr *string
	if latestDueDate.Valid {
		latestDueDatePtr = &latestDueDate.String
	}

	var borrowedDueDatePtr *string
	if currentLendingRow.DueDate.Valid {
		borrowedDueDatePtr = &currentLendingRow.DueDate.String
	}

	currentLending, err := domain.ParseCurrentLending(latestDueDatePtr, borrowedDueDatePtr, currentLendingRow.BorrowedAt)
	if err != nil {
		return nil, err
	}

	err = domain.ValidateReturn(currentLending, input.RequesterID, currentLendingRow.BorrowerID)
	if err != nil {
		if errors.Is(err, domain.ErrBookNotBorrowed) {
			return nil, ErrBookNotBorrowed
		}
		if errors.Is(err, domain.ErrNotBorrower) {
			return nil, ErrNotBorrower
		}
		return nil, err
	}

	eventID := uuid.New().String()
	err = s.lendingQueries.InsertLendingEvent(ctx, InsertLendingEventParams{
		EventID:    eventID,
		LendingID:  currentLendingRow.LendingID,
		BookID:     currentLendingRow.BookID,
		BorrowerID: currentLendingRow.BorrowerID,
		EventType:  "returned",
		DueDate:    sql.NullString{Valid: false},
		OccurredAt: now.Format(time.RFC3339),
	})
	if err != nil {
		return nil, err
	}

	return &ReturnBookOutput{
		LendingID:  currentLendingRow.LendingID,
		BookID:     currentLendingRow.BookID,
		BorrowerID: currentLendingRow.BorrowerID,
		ReturnedAt: now,
	}, nil
}
