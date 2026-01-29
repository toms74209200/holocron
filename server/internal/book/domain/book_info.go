package domain

import (
	"encoding/json"
	"errors"
	"strings"
)

var ErrBookNotFound = errors.New("book not found")

type BookInfo struct {
	Title         string
	Authors       []string
	Publisher     string
	PublishedDate string
	ThumbnailURL  string
}

func BookInfoFromGoogleBooks(body []byte) (*BookInfo, error) {
	var resp struct {
		TotalItems int `json:"totalItems"`
		Items      []struct {
			VolumeInfo struct {
				Title         string   `json:"title"`
				Authors       []string `json:"authors"`
				Publisher     string   `json:"publisher"`
				PublishedDate string   `json:"publishedDate"`
				ImageLinks    struct {
					Thumbnail string `json:"thumbnail"`
				} `json:"imageLinks"`
			} `json:"volumeInfo"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.TotalItems == 0 || len(resp.Items) == 0 {
		return nil, ErrBookNotFound
	}

	item := resp.Items[0].VolumeInfo
	return &BookInfo{
		Title:         item.Title,
		Authors:       item.Authors,
		Publisher:     item.Publisher,
		PublishedDate: item.PublishedDate,
		ThumbnailURL:  item.ImageLinks.Thumbnail,
	}, nil
}

func BookInfoFromOpenBD(body []byte) (*BookInfo, error) {
	var resp []struct {
		Summary *struct {
			Title     string `json:"title"`
			Author    string `json:"author"`
			Publisher string `json:"publisher"`
			Pubdate   string `json:"pubdate"`
			Cover     string `json:"cover"`
		} `json:"summary"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if len(resp) == 0 || resp[0].Summary == nil {
		return nil, ErrBookNotFound
	}

	summary := resp[0].Summary
	if summary.Title == "" {
		return nil, ErrBookNotFound
	}

	var authors []string
	if summary.Author != "" {
		authors = strings.Split(summary.Author, ",")
		for i := range authors {
			authors[i] = strings.TrimSpace(authors[i])
		}
	}

	return &BookInfo{
		Title:         summary.Title,
		Authors:       authors,
		Publisher:     summary.Publisher,
		PublishedDate: normalizeDate(summary.Pubdate),
		ThumbnailURL:  summary.Cover,
	}, nil
}

func normalizeDate(s string) string {
	if len(s) == 8 {
		return s[:4] + "-" + s[4:6] + "-" + s[6:8]
	}
	return s
}
