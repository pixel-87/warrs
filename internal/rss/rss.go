package rss

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/pixel-87/warss/internal/models"
	"github.com/pixel-87/warss/internal/storage"
)

// Allows for reuse of gofeed.Parser and http.client
type Fetcher struct {
	parser *gofeed.Parser
	client *http.Client
	db     *storage.DB
}

func NewFetcher(db *storage.DB) *Fetcher {
	return &Fetcher{
		parser: gofeed.NewParser(),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		db: db,
	}
}

func (f *Fetcher) fetchURL(url string) ([]byte, error) {
	resp, err := f.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("error closing response body %v", cerr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, nil
}

func (f *Fetcher) parseFeed(url string, data []byte) (models.Feed, error) {
	// Parses any feed into a universal gofeed.Feed, takes an io reader which reads xml/json data
	rawFeed, err := f.parser.Parse(strings.NewReader(string(data)))
	if err != nil {
		return models.Feed{}, fmt.Errorf("failed parsing %s: %w", url, err)
	}

	myFeed := models.Feed{
		Title: rawFeed.Title,
		URL:   url,
	}
	for _, item := range rawFeed.Items {
		content := item.Content
		if content == "" {
			content = item.Description
		}
		myFeed.Posts = append(myFeed.Posts, models.Post{
			Title:   item.Title,
			Link:    item.Link,
			Content: content,
		})
	}
	return myFeed, nil
}

func (f *Fetcher) GetFeed(url string) (models.Feed, error) {
	body, err := f.fetchURL(url)
	if err != nil {
		return models.Feed{}, err
	}

	return f.parseFeed(url, body)
}
