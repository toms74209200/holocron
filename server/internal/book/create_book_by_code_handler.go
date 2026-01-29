package book

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"holocron/internal/book/domain"
)

type CreateBookByCodeHandler struct {
	queries *Queries
	sources []domain.BookInfoSource
}

func NewCreateBookByCodeHandler(queries *Queries, sources []domain.BookInfoSource) *CreateBookByCodeHandler {
	return &CreateBookByCodeHandler{
		queries: queries,
		sources: sources,
	}
}

func (h *CreateBookByCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	output, err := CreateBookByCode(r.Context(), h.queries, h.sources, CreateBookByCodeInput{
		Code: req.Code,
	})

	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCode):
			writeError(w, http.StatusBadRequest, "invalid_request", "code must not be empty")
		case errors.Is(err, domain.ErrBookNotFound):
			writeError(w, http.StatusNotFound, "not_found", "book not found in external APIs")
		case errors.Is(err, ErrInvalidTitle):
			writeError(w, http.StatusBadRequest, "invalid_request", "external API returned invalid title")
		case errors.Is(err, ErrInvalidAuthors):
			writeError(w, http.StatusBadRequest, "invalid_request", "external API returned invalid authors")
		default:
			writeError(w, http.StatusInternalServerError, "internal_error", "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"id":            output.ID,
		"code":          output.Code,
		"title":         output.Title,
		"authors":       output.Authors,
		"publisher":     output.Publisher,
		"publishedDate": output.PublishedDate,
		"thumbnailUrl":  output.ThumbnailURL,
		"status":        output.Status,
		"createdAt":     output.CreatedAt.Format(time.RFC3339),
	})
}
