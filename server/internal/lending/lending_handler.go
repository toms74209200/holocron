package lending

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"holocron/internal/auth"
	"holocron/internal/book"
	"holocron/internal/lending/domain"
)

type BorrowBookHandler struct {
	service *BorrowBookService
}

func NewBorrowBookHandler(service *BorrowBookService) *BorrowBookHandler {
	return &BorrowBookHandler{
		service: service,
	}
}

func (h *BorrowBookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("bookId")
	if bookID == "" {
		writeError(w, http.StatusBadRequest, "invalid_request", "bookId is required")
		return
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "authentication required")
		return
	}

	var req struct {
		DueDays *int `json:"dueDays"`
	}
	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if err != io.EOF {
				writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
				return
			}
		}
	}

	output, err := h.service.BorrowBook(r.Context(), BorrowBookInput{
		BookID:     bookID,
		BorrowerID: userID,
		DueDays:    req.DueDays,
	})

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidDueDays):
			writeError(w, http.StatusBadRequest, "invalid_request", "due days must be at least 1")
		case errors.Is(err, ErrBookAlreadyBorrowed):
			writeError(w, http.StatusConflict, "book_already_borrowed", "book is already borrowed by another user")
		case errors.Is(err, ErrBookNotFound):
			writeError(w, http.StatusNotFound, "book_not_found", "book not found")
		default:
			writeError(w, http.StatusInternalServerError, "internal_error", "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"id":         output.ID,
		"bookId":     output.BookID,
		"borrowerId": output.BorrowerID,
		"borrowedAt": output.BorrowedAt.Format(time.RFC3339),
		"dueDate":    output.DueDate.Format(time.RFC3339),
	})
}

type ReturnBookHandler struct {
	service     *ReturnBookService
	bookQueries *book.Queries
}

func NewReturnBookHandler(service *ReturnBookService, bookQueries *book.Queries) *ReturnBookHandler {
	return &ReturnBookHandler{
		service:     service,
		bookQueries: bookQueries,
	}
}

func (h *ReturnBookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("bookId")
	if bookID == "" {
		writeError(w, http.StatusBadRequest, "invalid_request", "bookId is required")
		return
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "authentication required")
		return
	}

	_, err := h.service.ReturnBook(r.Context(), ReturnBookInput{
		BookID:      bookID,
		RequesterID: userID,
	})

	if err != nil {
		switch {
		case errors.Is(err, ErrBookNotBorrowed):
			writeError(w, http.StatusConflict, "not_borrowed", "this book is not currently borrowed")
		case errors.Is(err, ErrNotBorrower):
			writeError(w, http.StatusForbidden, "forbidden", "only the borrower can return this book")
		case errors.Is(err, ErrBookNotFound):
			writeError(w, http.StatusNotFound, "book_not_found", "book not found")
		default:
			writeError(w, http.StatusInternalServerError, "internal_error", "internal server error")
		}
		return
	}

	bookDetails, err := book.GetBook(r.Context(), h.bookQueries, book.GetBookInput{
		BookID: bookID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "failed to retrieve book details")
		return
	}

	resp := map[string]any{
		"id":        bookDetails.ID,
		"title":     bookDetails.Title,
		"authors":   bookDetails.Authors,
		"status":    "available",
		"createdAt": bookDetails.CreatedAt.Format(time.RFC3339),
	}
	if bookDetails.Code != nil {
		resp["code"] = *bookDetails.Code
	}
	if bookDetails.Publisher != nil {
		resp["publisher"] = *bookDetails.Publisher
	}
	if bookDetails.PublishedDate != nil {
		resp["publishedDate"] = *bookDetails.PublishedDate
	}
	if bookDetails.ThumbnailURL != nil {
		resp["thumbnailUrl"] = *bookDetails.ThumbnailURL
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"code":    code,
		"message": message,
	})
}
