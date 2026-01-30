package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"holocron/internal/api"
	"holocron/internal/auth"
	"holocron/internal/book"
	"holocron/internal/book/domain"
	"holocron/internal/user"

	_ "modernc.org/sqlite"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

type server struct {
	createUserHandler       *user.CreateUserHandler
	createBookHandler       *book.CreateBookHandler
	createBookByCodeHandler *book.CreateBookByCodeHandler
	listBooksHandler        *book.ListBooksHandler
}

func (s *server) GetBooks(w http.ResponseWriter, r *http.Request, params api.GetBooksParams) {
	s.listBooksHandler.ServeHTTP(w, r, params)
}
func (s *server) PostBooks(w http.ResponseWriter, r *http.Request) {
	s.createBookHandler.ServeHTTP(w, r)
}
func (s *server) PostBooksCode(w http.ResponseWriter, r *http.Request) {
	s.createBookByCodeHandler.ServeHTTP(w, r)
}
func (s *server) GetBook(w http.ResponseWriter, r *http.Request, bookId openapi_types.UUID) {
	notImplemented(w)
}
func (s *server) PostBooksBookId(w http.ResponseWriter, r *http.Request, bookId openapi_types.UUID) {
	notImplemented(w)
}
func (s *server) PostBooksBorrow(w http.ResponseWriter, r *http.Request, bookId openapi_types.UUID) {
	notImplemented(w)
}
func (s *server) PostBooksReturn(w http.ResponseWriter, r *http.Request, bookId openapi_types.UUID) {
	notImplemented(w)
}

func (s *server) PostUsers(w http.ResponseWriter, r *http.Request) {
	s.createUserHandler.ServeHTTP(w, r)
}

func notImplemented(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "not implemented"})
}

func initDB(database *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS user_events (
		event_id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		event_type TEXT NOT NULL,
		name TEXT NOT NULL,
		occurred_at TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_user_events_user_id ON user_events(user_id);

	CREATE TABLE IF NOT EXISTS book_events (
		event_id TEXT PRIMARY KEY,
		book_id TEXT NOT NULL,
		event_type TEXT NOT NULL,
		code TEXT,
		title TEXT,
		authors TEXT,
		publisher TEXT,
		published_date TEXT,
		thumbnail_url TEXT,
		occurred_at TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_book_events_book_id ON book_events(book_id);
	`
	_, err := database.Exec(schema)
	return err
}

func main() {
	ctx := context.Background()

	database, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	if err := initDB(database); err != nil {
		log.Fatal(err)
	}

	firebaseAuth, err := auth.NewFirebaseAuth(ctx)
	if err != nil {
		log.Fatal(err)
	}

	userQueries := user.New(database)
	bookQueries := book.New(database)

	googleBooksFetcher, err := book.NewGoogleBooksFetcher()
	if err != nil {
		log.Fatal(err)
	}
	openBDFetcher, err := book.NewOpenBDFetcher()
	if err != nil {
		log.Fatal(err)
	}

	bookInfoSources := []domain.BookInfoSource{
		book.DBCacheSource(bookQueries),
		book.ExternalAPISource(googleBooksFetcher.Fetch, domain.BookInfoFromGoogleBooks),
		book.ExternalAPISource(openBDFetcher.Fetch, domain.BookInfoFromOpenBD),
	}

	srv := &server{
		createUserHandler:       user.NewCreateUserHandler(userQueries, firebaseAuth),
		createBookHandler:       book.NewCreateBookHandler(bookQueries),
		createBookByCodeHandler: book.NewCreateBookByCodeHandler(bookQueries, bookInfoSources),
		listBooksHandler:        book.NewListBooksHandler(bookQueries),
	}

	mux := http.NewServeMux()
	api.HandlerWithOptions(srv, api.StdHTTPServerOptions{
		BaseRouter: mux,
		Middlewares: []api.MiddlewareFunc{
			auth.AuthMiddleware(firebaseAuth),
		},
	})

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
