package main

import (
	"fmt"

	"github.com/pixel-87/edwarss/internal/rss"
)

func main() {
	fetcher := rss.NewFetcher()
	url := "https://ed-thomas.dev/rss.xml"

	feed, err := fetcher.GetFeed(url)
	if err != nil {
		fmt.Printf("err %v", err)
		return
	}
	fmt.Println(feed.Title)
}
