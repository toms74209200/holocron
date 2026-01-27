package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"holocron/internal/api"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

type stub struct{}

func (s *stub) CreateBook(w http.ResponseWriter, r *http.Request)                                   { notImplemented(w) }
func (s *stub) ListBooks(w http.ResponseWriter, r *http.Request, params api.ListBooksParams)        { notImplemented(w) }
func (s *stub) CreateBookByCode(w http.ResponseWriter, r *http.Request)                             { notImplemented(w) }
func (s *stub) GetBook(w http.ResponseWriter, r *http.Request, bookId openapi_types.UUID)           { notImplemented(w) }
func (s *stub) UpdateBook(w http.ResponseWriter, r *http.Request, bookId openapi_types.UUID)        { notImplemented(w) }
func (s *stub) BorrowBook(w http.ResponseWriter, r *http.Request, bookId openapi_types.UUID)        { notImplemented(w) }
func (s *stub) ReturnBook(w http.ResponseWriter, r *http.Request, bookId openapi_types.UUID)        { notImplemented(w) }
func (s *stub) CreateUser(w http.ResponseWriter, r *http.Request)                                   { notImplemented(w) }

func notImplemented(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{"message": "not implemented"})
}

func main() {
	mux := http.NewServeMux()
	api.HandlerFromMux(&stub{}, mux)

	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", mux)
}
