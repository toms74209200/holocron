package book

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"holocron/internal/book/domain"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

type GetBookHandler struct {
	queries *Queries
}

func NewGetBookHandler(queries *Queries) *GetBookHandler {
	return &GetBookHandler{
		queries: queries,
	}
}

func (h *GetBookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, bookId openapi_types.UUID) {
	output, err := GetBook(r.Context(), h.queries, GetBookInput{
		BookID: bookId.String(),
	})

	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			writeError(w, http.StatusNotFound, "not_found", "book not found")
			return
		}
		if errors.Is(err, ErrInvalidBookID) {
			writeError(w, http.StatusBadRequest, "invalid_request", "invalid book ID")
			return
		}
		if errors.Is(err, ErrInvalidBookRow) {
			writeError(w, http.StatusInternalServerError, "internal_error", "invalid book data")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "internal server error")
		return
	}

	resp := map[string]any{
		"id":        output.ID,
		"title":     output.Title,
		"authors":   output.Authors,
		"status":    output.Status,
		"createdAt": output.CreatedAt.Format(time.RFC3339),
	}
	if output.Code != nil {
		resp["code"] = *output.Code
	}
	if output.Publisher != nil {
		resp["publisher"] = *output.Publisher
	}
	if output.PublishedDate != nil {
		resp["publishedDate"] = *output.PublishedDate
	}
	if output.ThumbnailURL != nil {
		resp["thumbnailUrl"] = *output.ThumbnailURL
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

type UpdateBookHandler struct {
	queries *Queries
}

func NewUpdateBookHandler(queries *Queries) *UpdateBookHandler {
	return &UpdateBookHandler{
		queries: queries,
	}
}

func (h *UpdateBookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, bookId openapi_types.UUID) {
	var req struct {
		Code          *string   `json:"code"`
		Title         *string   `json:"title"`
		Authors       *[]string `json:"authors"`
		Publisher     *string   `json:"publisher"`
		PublishedDate *string   `json:"publishedDate"`
		ThumbnailURL  *string   `json:"thumbnailUrl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	output, err := UpdateBook(r.Context(), h.queries, UpdateBookInput{
		BookID:        bookId.String(),
		Code:          req.Code,
		Title:         req.Title,
		Authors:       req.Authors,
		Publisher:     req.Publisher,
		PublishedDate: req.PublishedDate,
		ThumbnailURL:  req.ThumbnailURL,
	})

	if err != nil {
		switch {
		case errors.Is(err, ErrBookNotFound):
			writeError(w, http.StatusNotFound, "not_found", "book not found")
		case errors.Is(err, domain.ErrInvalidTitle):
			writeError(w, http.StatusBadRequest, "invalid_request", "title must be 1-200 characters")
		case errors.Is(err, domain.ErrInvalidAuthors):
			writeError(w, http.StatusBadRequest, "invalid_request", "authors must have at least one author")
		case errors.Is(err, ErrInvalidBookRow):
			writeError(w, http.StatusInternalServerError, "internal_error", "invalid book data")
		default:
			writeError(w, http.StatusInternalServerError, "internal_error", "internal server error")
		}
		return
	}

	resp := map[string]any{
		"id":        output.ID,
		"title":     output.Title,
		"authors":   output.Authors,
		"status":    output.Status,
		"createdAt": output.CreatedAt.Format(time.RFC3339),
		"updatedAt": output.UpdatedAt.Format(time.RFC3339),
	}
	if output.Code != nil {
		resp["code"] = *output.Code
	}
	if output.Publisher != nil {
		resp["publisher"] = *output.Publisher
	}
	if output.PublishedDate != nil {
		resp["publishedDate"] = *output.PublishedDate
	}
	if output.ThumbnailURL != nil {
		resp["thumbnailUrl"] = *output.ThumbnailURL
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
