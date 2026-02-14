package bookcode

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type OpenBDFetcher struct {
	baseURL string
	client  *http.Client
}

func NewOpenBDFetcher() (*OpenBDFetcher, error) {
	baseURL := os.Getenv("OPENBD_API_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("OPENBD_API_URL is not set")
	}
	return &OpenBDFetcher{
		baseURL: baseURL,
		client: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}, nil
}

func (f *OpenBDFetcher) Fetch(ctx context.Context, code string) ([]byte, error) {
	reqURL := fmt.Sprintf("%s/get?isbn=%s", f.baseURL, url.QueryEscape(code))

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
		return nil, fmt.Errorf("openbd api returned status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
