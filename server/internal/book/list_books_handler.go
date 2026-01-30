package book

import (
	"encoding/json"
	"net/http"
	"time"

	"holocron/internal/api"
)

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
		respItems = append(respItems, m)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"items": respItems,
		"total": output.Total,
	})
}
