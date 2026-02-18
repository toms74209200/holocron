package books

import (
	"encoding/json"
	"errors"
	"holocron/internal/api"
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
		Code          *string  `json:"code"`
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
		Code:          req.Code,
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

type ListBooksHandler struct {
	queries *Queries
}

func NewListBooksHandler(queries *Queries) *ListBooksHandler {
	return &ListBooksHandler{
		queries: queries,
	}
}

func (h *ListBooksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, params api.GetBooksParams) {
	output, err := ListBooks(r.Context(), h.queries, ListBooksInput{
		Q:      params.Q,
		Code:   params.Code,
		Limit:  params.Limit,
		Offset: params.Offset,
	})

	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "internal server error")
		return
	}

	respItems := make([]map[string]any, 0, len(output.Items))
	for _, item := range output.Items {
		m := map[string]any{
			"id":        item.ID,
			"title":     item.Title,
			"authors":   item.Authors,
			"status":    item.Status,
			"createdAt": item.CreatedAt.Format(time.RFC3339),
		}
		if item.Code != nil {
			m["code"] = *item.Code
		}
		if item.Publisher != nil {
			m["publisher"] = *item.Publisher
		}
		if item.PublishedDate != nil {
			m["publishedDate"] = *item.PublishedDate
		}
		if item.ThumbnailURL != nil {
			m["thumbnailUrl"] = *item.ThumbnailURL
		}
		if item.Borrower != nil {
			m["borrower"] = map[string]any{
				"id":         item.Borrower.ID,
				"name":       item.Borrower.Name,
				"borrowedAt": item.Borrower.BorrowedAt.Format(time.RFC3339),
			}
		}
		respItems = append(respItems, m)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"items": respItems,
		"total": output.Total,
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
