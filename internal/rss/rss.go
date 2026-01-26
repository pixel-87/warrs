package rss

import "github.com/mmcdole/gofeed"

// A single post in a feed
type post struct {
	title string
	content string
	link string
	id int
	read bool
}

// An entire feed
type feed struct {
	title string
	URL string
	posts []post
	id int
}

func GetFeed(url string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	return fp.ParseURL(url)
}
