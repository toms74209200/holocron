package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"holocron/internal/auth"
	"holocron/internal/lending"
)

type CreateUserHandler struct {
	queries      *Queries
	firebaseAuth FirebaseAuth
}

func NewCreateUserHandler(queries *Queries, firebaseAuth FirebaseAuth) *CreateUserHandler {
	return &CreateUserHandler{
		queries:      queries,
		firebaseAuth: firebaseAuth,
	}
}

func (h *CreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name *string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	output, err := CreateUser(r.Context(), h.queries, h.firebaseAuth, CreateUserInput{
		Name: req.Name,
	})

	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidUserName):
			writeError(w, http.StatusBadRequest, "invalid_request", "name must be 1-50 characters")
		case errors.Is(err, ErrUserAlreadyExists):
			writeError(w, http.StatusConflict, "user_exists", "user already exists")
		case errors.Is(err, ErrTokenCreation):
			writeError(w, http.StatusInternalServerError, "internal_error", "failed to create token")
		default:
			writeError(w, http.StatusInternalServerError, "internal_error", "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          output.ID,
		"name":        output.Name,
		"customToken": output.CustomToken,
		"createdAt":   output.CreatedAt.Format(time.RFC3339),
	})
}

type GetMyBorrowingHandler struct {
	lendingQueries *lending.Queries
}

func NewGetMyBorrowingHandler(lendingQueries *lending.Queries) *GetMyBorrowingHandler {
	return &GetMyBorrowingHandler{
		lendingQueries: lendingQueries,
	}
}

func (h *GetMyBorrowingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "authentication required")
		return
	}

	output, err := GetMyBorrowing(r.Context(), h.lendingQueries, GetMyBorrowingInput{
		BorrowerID: userID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "internal server error")
		return
	}

	items := make([]map[string]any, 0, len(output.Items))
	for _, item := range output.Items {
		m := map[string]any{
			"id":         item.ID,
			"title":      item.Title,
			"authors":    item.Authors,
			"borrowedAt": item.BorrowedAt.Format(time.RFC3339),
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
		if item.DueDate != nil {
			m["dueDate"] = item.DueDate.Format(time.RFC3339)
		}
		items = append(items, m)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"items": items,
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
