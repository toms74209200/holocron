package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"holocron/internal/api"
	"holocron/internal/auth"
	"holocron/internal/book"
	bookDomain "holocron/internal/book/domain"
	"holocron/internal/bookcode"
	bookcodeDomain "holocron/internal/bookcode/domain"
	"holocron/internal/books"
	"holocron/internal/lending"
	"holocron/internal/user"

	_ "modernc.org/sqlite"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

type server struct {
	createUserHandler       *user.CreateUserHandler
	createBookHandler       *books.CreateBookHandler
	createBookByCodeHandler *bookcode.CreateBookByCodeHandler
	listBooksHandler        *books.ListBooksHandler
	getBookHandler          *book.GetBookHandler
	updateBookHandler       *book.UpdateBookHandler
	deleteBookHandler       *book.DeleteBookHandler
	borrowBookHandler       *lending.BorrowBookHandler
	returnBookHandler       *lending.ReturnBookHandler
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
	s.getBookHandler.ServeHTTP(w, r, bookId)
}
func (s *server) PostBooksBookId(w http.ResponseWriter, r *http.Request, bookId openapi_types.UUID) {
	s.updateBookHandler.ServeHTTP(w, r, bookId)
}
func (s *server) DeleteBook(w http.ResponseWriter, r *http.Request, bookId openapi_types.UUID) {
	s.deleteBookHandler.ServeHTTP(w, r, bookId)
}
func (s *server) PostBooksBorrow(w http.ResponseWriter, r *http.Request, bookId openapi_types.UUID) {
	s.borrowBookHandler.ServeHTTP(w, r)
}
func (s *server) PostBooksReturn(w http.ResponseWriter, r *http.Request, bookId openapi_types.UUID) {
	s.returnBookHandler.ServeHTTP(w, r)
}

func (s *server) PostUsers(w http.ResponseWriter, r *http.Request) {
	s.createUserHandler.ServeHTTP(w, r)
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
		delete_reason TEXT,
		delete_memo TEXT,
		occurred_at TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_book_events_book_id ON book_events(book_id);

	CREATE TABLE IF NOT EXISTS lending_events (
		event_id TEXT PRIMARY KEY,
		lending_id TEXT NOT NULL,
		book_id TEXT NOT NULL,
		borrower_id TEXT NOT NULL,
		event_type TEXT NOT NULL,
		due_date TEXT,
		occurred_at TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_lending_events_lending_id ON lending_events(lending_id);
	CREATE INDEX IF NOT EXISTS idx_lending_events_book_id ON lending_events(book_id);
	CREATE INDEX IF NOT EXISTS idx_lending_events_borrower_id ON lending_events(borrower_id);
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
	booksQueries := books.New(database)
	bookcodeQueries := bookcode.New(database)
	bookQueries := book.New(database)
	lendingQueries := lending.New(database)

	googleBooksFetcher, err := bookcode.NewGoogleBooksFetcher()
	if err != nil {
		log.Fatal(err)
	}
	openBDFetcher, err := bookcode.NewOpenBDFetcher()
	if err != nil {
		log.Fatal(err)
	}

	bookInfoSources := []bookcodeDomain.BookInfoSource{
		bookcode.DBCacheSource(bookcodeQueries),
		bookcode.ExternalAPISource(googleBooksFetcher.Fetch, bookDomain.BookInfoFromGoogleBooks),
		bookcode.ExternalAPISource(openBDFetcher.Fetch, bookDomain.BookInfoFromOpenBD),
	}

	borrowBookService := lending.NewBorrowBookService(lendingQueries, bookQueries)
	returnBookService := lending.NewReturnBookService(lendingQueries, bookQueries)

	srv := &server{
		createUserHandler:       user.NewCreateUserHandler(userQueries, firebaseAuth),
		createBookHandler:       books.NewCreateBookHandler(booksQueries),
		createBookByCodeHandler: bookcode.NewCreateBookByCodeHandler(bookcodeQueries, bookInfoSources),
		listBooksHandler:        books.NewListBooksHandler(booksQueries),
		getBookHandler:          book.NewGetBookHandler(bookQueries),
		updateBookHandler:       book.NewUpdateBookHandler(bookQueries),
		deleteBookHandler:       book.NewDeleteBookHandler(bookQueries),
		borrowBookHandler:       lending.NewBorrowBookHandler(borrowBookService),
		returnBookHandler:       lending.NewReturnBookHandler(returnBookService, bookQueries),
	}

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	api.HandlerWithOptions(srv, api.StdHTTPServerOptions{
		BaseRouter: mux,
		Middlewares: []api.MiddlewareFunc{
			auth.CORSMiddleware(allowedOrigin),
			auth.AuthMiddleware(firebaseAuth),
		},
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		mux.ServeHTTP(w, r)
	})

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	serverErrors := make(chan error, 1)
	go func() {
		fmt.Println("Server starting on :8080")
		serverErrors <- httpServer.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	case sig := <-signalChan:
		log.Printf("Received signal: %v. Starting graceful shutdown...", sig)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("Error during server shutdown: %v", err)
			if err := httpServer.Close(); err != nil {
				log.Printf("Error closing server: %v", err)
			}
		}

		log.Println("Server stopped gracefully")
	}
}
