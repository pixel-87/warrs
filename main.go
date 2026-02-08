package main

import (
	"fmt"
	"log"

	"github.com/pixel-87/warss/internal/rss"
	"github.com/pixel-87/warss/internal/storage"
)

func main() {
	db, err := storage.NewDB("./rss.db")
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
	defer db.Close()

	fetcher := rss.NewFetcher(db)

	// 1. Ensure our test feed is in there
	_ = db.AddFeed("https://ed-thomas.dev/rss.xml", "Ed's Blog")

	// 2. Get everything we are subscribed to
	subscriptions, err := db.GetFeeds()
	if err != nil {
		log.Fatalf("could not get subscriptions: %v", err)
	}

	// 3. Process the collection
	fmt.Printf("Updating %d feeds...\n", len(subscriptions))
	for _, s := range subscriptions {
		feed, err := fetcher.GetFeed(s.URL) // Using the URL from the DB!
		if err != nil {
			fmt.Printf("Failed to update %s: %v\n", s.Title, err)
			continue
		}
		fmt.Printf("✅ %s (%d posts)\n", feed.Title, len(feed.Posts))
	}
}
