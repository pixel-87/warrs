package main

import (
	"fmt"

	"github.com/pixel-87/edwarss/internal/rss"
)

func main() {
	url := "https://ed-thomas.dev/rss.xml"
	feed, err := rss.GetFeed(url)
	if err != nil {
		fmt.Printf("err %v", err)
		return 
	}
	fmt.Println(feed.Title)
}
