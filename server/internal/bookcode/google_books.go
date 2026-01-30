package bookcode

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type GoogleBooksFetcher struct {
	baseURL string
	client  *http.Client
}

func NewGoogleBooksFetcher() (*GoogleBooksFetcher, error) {
	baseURL := os.Getenv("GOOGLE_BOOKS_API_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("GOOGLE_BOOKS_API_URL is not set")
	}
	return &GoogleBooksFetcher{
		baseURL: baseURL,
		client:  http.DefaultClient,
	}, nil
}

func (f *GoogleBooksFetcher) Fetch(ctx context.Context, code string) ([]byte, error) {
	reqURL := fmt.Sprintf("%s/volumes?q=%s", f.baseURL, url.QueryEscape("isbn:"+code))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google books api returned status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
