package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type CreateUserHandler struct {
	service *CreateUserService
}

func NewCreateUserHandler(service *CreateUserService) *CreateUserHandler {
	return &CreateUserHandler{
		service: service,
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

	output, err := h.service.Execute(r.Context(), CreateUserInput{
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

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"code":    code,
		"message": message,
	})
}
