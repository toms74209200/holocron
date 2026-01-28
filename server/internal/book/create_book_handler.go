package book

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type CreateBookHandler struct {
	queries *Queries
}

func NewCreateBookHandler(queries *Queries) *CreateBookHandler {
	return &CreateBookHandler{
		queries: queries,
	}
}

func (h *CreateBookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title         string   `json:"title"`
		Authors       []string `json:"authors"`
		Publisher     *string  `json:"publisher"`
		PublishedDate *string  `json:"publishedDate"`
		ThumbnailURL  *string  `json:"thumbnailUrl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	output, err := CreateBook(r.Context(), h.queries, CreateBookInput{
		Title:         req.Title,
		Authors:       req.Authors,
		Publisher:     req.Publisher,
		PublishedDate: req.PublishedDate,
		ThumbnailURL:  req.ThumbnailURL,
	})

	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidTitle):
			writeError(w, http.StatusBadRequest, "invalid_request", "title must be 1-200 characters")
		case errors.Is(err, ErrInvalidAuthors):
			writeError(w, http.StatusBadRequest, "invalid_request", "authors must have at least one author")
		default:
			writeError(w, http.StatusInternalServerError, "internal_error", "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"id":            output.ID,
		"title":         output.Title,
		"authors":       output.Authors,
		"publisher":     output.Publisher,
		"publishedDate": output.PublishedDate,
		"thumbnailUrl":  output.ThumbnailURL,
		"status":        output.Status,
		"createdAt":     output.CreatedAt.Format(time.RFC3339),
	})
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"code":    code,
		"message": message,
	})
}
