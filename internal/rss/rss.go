package rss

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

// A single post in a feed
type Post struct {
	Title   string
	Content string
	Link    string
	ID      int
	Read    bool
}

// An entire Feed
type Feed struct {
	Title string
	URL   string
	Posts []Post
	ID    int
}

type Fetcher struct {
	parser *gofeed.Parser
	client *http.Client
}

func NewFetcher() *Fetcher {
	return &Fetcher{
		parser: gofeed.NewParser(),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (f *Fetcher) fetchURL(url string) ([]byte, error) {
	resp, err := f.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer func () {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("error closing response body %v", cerr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read responce body: %w", err)
	}
	return body, nil
}

func (f *Fetcher) parseFeed(url string, data []byte) (Feed, error) {
	rawFeed, err := f.parser.Parse(strings.NewReader(string(data)))
	if err != nil {
		return Feed{}, fmt.Errorf("failed parsing %s: %w", url, err)
	}

	myFeed := Feed{
		Title: rawFeed.Title,
		URL:   url,
	}
	for _, item := range rawFeed.Items {
		content := item.Content
		if content == "" {
			content = item.Description
		}
		myFeed.Posts = append(myFeed.Posts, Post{
			Title:   item.Title,
			Link:    item.Link,
			Content: content,
		})
	}
	return myFeed, nil
}

func (f *Fetcher) GetFeed(url string) (Feed, error) {
	body, err := f.fetchURL(url)
	if err != nil {
		return Feed{}, err
	}

	return f.parseFeed(url, body)
}
