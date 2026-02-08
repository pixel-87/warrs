package storage

import (
	"os"
	"testing"

	"github.com/pixel-87/warss/internal/models"
)

// Helper to create a temporary test database
func setupTestDB(t *testing.T) *DB {
	tmpFile := "test_feeds.db"
	// Clean up any previous test file
	os.Remove(tmpFile)

	db, err := NewDB(tmpFile)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
		os.Remove(tmpFile)
	})

	return db
}

func TestAddFeed(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name    string
		url     string
		title   string
		wantErr bool
	}{
		{
			name:    "Add valid feed",
			url:     "https://example.com/feed.xml",
			title:   "Example Feed",
			wantErr: false,
		},
		{
			name:    "Add feed with empty title",
			url:     "https://example.com/feed2.xml",
			title:   "",
			wantErr: false,
		},
		{
			name:    "Add duplicate URL",
			url:     "https://example.com/feed.xml",
			title:   "Duplicate",
			wantErr: true, // UNIQUE constraint
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.AddFeed(tt.url, tt.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddFeed() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetFeeds(t *testing.T) {
	db := setupTestDB(t)

	// Add some test feeds
	feeds := []struct {
		url   string
		title string
	}{
		{"https://example.com/feed1.xml", "Feed 1"},
		{"https://example.com/feed2.xml", "Feed 2"},
		{"https://example.com/feed3.xml", "Feed 3"},
	}

	for _, feed := range feeds {
		if err := db.AddFeed(feed.url, feed.title); err != nil {
			t.Fatalf("failed to add test feed: %v", err)
		}
	}

	// Get feeds
	result, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("GetFeeds() error = %v", err)
	}

	if len(result) != len(feeds) {
		t.Errorf("GetFeeds() got %d feeds, want %d", len(result), len(feeds))
	}

	// Verify fields are populated correctly
	for i, feed := range result {
		if feed.ID == 0 {
			t.Errorf("Feed %d has zero ID", i)
		}
		if feed.URL != feeds[i].url {
			t.Errorf("Feed %d URL = %q, want %q", i, feed.URL, feeds[i].url)
		}
		if feed.Title != feeds[i].title {
			t.Errorf("Feed %d Title = %q, want %q", i, feed.Title, feeds[i].title)
		}
	}
}

func TestUpdateFeed(t *testing.T) {
	db := setupTestDB(t)

	// Add a feed
	if err := db.AddFeed("https://example.com/old.xml", "Old Title"); err != nil {
		t.Fatalf("failed to add test feed: %v", err)
	}

	// Get it to get the ID
	feeds, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("GetFeeds() error = %v", err)
	}

	if len(feeds) == 0 {
		t.Fatalf("no feeds found")
	}

	// Update the feed
	feeds[0].URL = "https://example.com/new.xml"
	feeds[0].Title = "New Title"

	if err := db.UpdateFeed(feeds[0]); err != nil {
		t.Fatalf("UpdateFeed() error = %v", err)
	}

	// Verify the update
	updated, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("GetFeeds() error after update = %v", err)
	}

	if len(updated) != 1 {
		t.Fatalf("expected 1 feed, got %d", len(updated))
	}

	if updated[0].URL != "https://example.com/new.xml" {
		t.Errorf("URL not updated: got %q, want %q", updated[0].URL, "https://example.com/new.xml")
	}

	if updated[0].Title != "New Title" {
		t.Errorf("Title not updated: got %q, want %q", updated[0].Title, "New Title")
	}
}

func TestDeleteFeed(t *testing.T) {
	db := setupTestDB(t)

	// Add a feed
	if err := db.AddFeed("https://example.com/delete.xml", "To Delete"); err != nil {
		t.Fatalf("failed to add test feed: %v", err)
	}

	// Get it to get the ID
	feeds, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("GetFeeds() error = %v", err)
	}

	feedID := feeds[0].ID

	// Delete the feed
	if err := db.DeleteFeed(feedID); err != nil {
		t.Fatalf("DeleteFeed() error = %v", err)
	}

	// Verify it's deleted
	remaining, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("GetFeeds() error after delete = %v", err)
	}

	if len(remaining) != 0 {
		t.Errorf("expected 0 feeds after delete, got %d", len(remaining))
	}
}

func TestGetFeedsEmpty(t *testing.T) {
	db := setupTestDB(t)

	feeds, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("GetFeeds() error = %v", err)
	}

	if len(feeds) != 0 {
		t.Errorf("empty database should return 0 feeds, got %d", len(feeds))
	}
}

func TestMultipleFeedsOperations(t *testing.T) {
	db := setupTestDB(t)

	// Add multiple feeds
	feedData := []models.Feed{
		{URL: "https://blog1.com/feed.xml", Title: "Blog 1"},
		{URL: "https://blog2.com/feed.xml", Title: "Blog 2"},
		{URL: "https://blog3.com/feed.xml", Title: "Blog 3"},
	}

	for _, f := range feedData {
		if err := db.AddFeed(f.URL, f.Title); err != nil {
			t.Fatalf("failed to add feed: %v", err)
		}
	}

	// Get all feeds
	feeds, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("GetFeeds() error = %v", err)
	}

	initialCount := len(feeds)
	if initialCount != 3 {
		t.Errorf("expected 3 feeds, got %d", initialCount)
	}

	// Delete one feed
	if err := db.DeleteFeed(feeds[1].ID); err != nil {
		t.Fatalf("DeleteFeed() error = %v", err)
	}

	// Verify count decreased
	feeds, err = db.GetFeeds()
	if err != nil {
		t.Fatalf("GetFeeds() error = %v", err)
	}

	if len(feeds) != 2 {
		t.Errorf("expected 2 feeds after deletion, got %d", len(feeds))
	}
}
