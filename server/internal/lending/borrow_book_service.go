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
	ErrBookAlreadyBorrowed = errors.New("book is already borrowed by another user")
	ErrBookNotFound        = errors.New("book not found")
)

type BorrowBookInput struct {
	BookID     string
	BorrowerID string
	DueDays    *int
}

type BorrowBookOutput struct {
	ID         string
	BookID     string
	BorrowerID string
	BorrowedAt time.Time
	DueDate    time.Time
}

type BookQueries interface {
	CountBookByBookId(ctx context.Context, bookID string) (int64, error)
}

type BorrowBookService struct {
	lendingQueries *Queries
	bookQueries BookQueries
	now         func() time.Time
}

func NewBorrowBookService(lendingQueries *Queries, bookQueries BookQueries) *BorrowBookService {
	return &BorrowBookService{
		lendingQueries: lendingQueries,
		bookQueries: bookQueries,
		now:         func() time.Time { return time.Now().UTC() },
	}
}

func (s *BorrowBookService) BorrowBook(ctx context.Context, input BorrowBookInput) (*BorrowBookOutput, error) {
	now := s.now()

	count, err := s.bookQueries.CountBookByBookId(ctx, input.BookID)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, ErrBookNotFound
	}

	currentLendingRow, err := s.lendingQueries.GetCurrentLending(ctx, input.BookID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	var currentLending *domain.CurrentLending
	if !errors.Is(err, sql.ErrNoRows) {
		if currentLendingRow.BorrowerID != input.BorrowerID {
			return nil, ErrBookAlreadyBorrowed
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

		currentLending, err = domain.ParseCurrentLending(latestDueDatePtr, borrowedDueDatePtr, currentLendingRow.BorrowedAt)
		if err != nil {
			return nil, err
		}
	}

	dueDate, _, err := domain.CalculateDueDate(now, input.DueDays, currentLending)
	if err != nil {
		return nil, err
	}

	if currentLending == nil {
		lendingID := uuid.New().String()
		eventID := uuid.New().String()

		err := s.lendingQueries.InsertLendingEvent(ctx, InsertLendingEventParams{
			EventID:    eventID,
			LendingID:  lendingID,
			BookID:     input.BookID,
			BorrowerID: input.BorrowerID,
			EventType:  "borrowed",
			DueDate:    sql.NullString{String: dueDate.Format(time.RFC3339), Valid: true},
			OccurredAt: now.Format(time.RFC3339),
		})
		if err != nil {
			return nil, err
		}

		return &BorrowBookOutput{
			ID:         lendingID,
			BookID:     input.BookID,
			BorrowerID: input.BorrowerID,
			BorrowedAt: now,
			DueDate:    dueDate,
		}, nil
	}

	eventID := uuid.New().String()
	err = s.lendingQueries.InsertLendingEvent(ctx, InsertLendingEventParams{
		EventID:    eventID,
		LendingID:     currentLendingRow.LendingID,
		BookID:     currentLendingRow.BookID,
		BorrowerID: currentLendingRow.BorrowerID,
		EventType:  "due_date_extended",
		DueDate:    sql.NullString{String: dueDate.Format(time.RFC3339), Valid: true},
		OccurredAt: now.Format(time.RFC3339),
	})
	if err != nil {
		return nil, err
	}

	borrowedAt, err := time.Parse(time.RFC3339, currentLendingRow.BorrowedAt)
	if err != nil {
		return nil, err
	}

	return &BorrowBookOutput{
		ID:         currentLendingRow.LendingID,
		BookID:     currentLendingRow.BookID,
		BorrowerID: currentLendingRow.BorrowerID,
		BorrowedAt: borrowedAt,
		DueDate:    dueDate,
	}, nil
}
