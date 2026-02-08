package storage

import (
	"fmt"
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

// TestAddFeedEdgeCases tests edge cases with unusual input
func TestAddFeedEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		title   string
		wantErr bool
	}{
		{
			name:    "Very long URL",
			url:     "https://example.com/" + string(make([]byte, 2000)),
			title:   "Test",
			wantErr: false,
		},
		{
			name:    "Unicode in title",
			url:     "https://example.com/feed1.xml",
			title:   "日本語のタイトル 🚀",
			wantErr: false,
		},
		{
			name:    "Special characters in URL",
			url:     "https://example.com/feed?param=value&other=test",
			title:   "Query Params",
			wantErr: false,
		},
		{
			name:    "SQL-like content in title (injection test)",
			url:     "https://example.com/sql.xml",
			title:   "'; DROP TABLE feeds; --",
			wantErr: false,
		},
		{
			name:    "SQL-like content in URL (injection test)",
			url:     "https://example.com/feed'; DELETE FROM feeds WHERE '1'='1",
			title:   "Injection Test",
			wantErr: false,
		},
		{
			name:    "Newlines in title",
			url:     "https://example.com/newline.xml",
			title:   "Line 1\nLine 2\nLine 3",
			wantErr: false,
		},
		{
			name:    "Tabs and spaces in title",
			url:     "https://example.com/tabs.xml",
			title:   "Title\t\twith\ttabs\tand   spaces",
			wantErr: false,
		},
		{
			name:    "HTML in title",
			url:     "https://example.com/html.xml",
			title:   "<script>alert('xss')</script>",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			
			err := db.AddFeed(tt.url, tt.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddFeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify it was stored correctly if no error expected
			if !tt.wantErr {
				feeds, err := db.GetFeeds()
				if err != nil {
					t.Fatalf("GetFeeds() error = %v", err)
				}
				
				if len(feeds) != 1 {
					t.Fatalf("expected 1 feed, got %d", len(feeds))
				}
				
				if feeds[0].URL != tt.url {
					t.Errorf("URL mismatch: got %q, want %q", feeds[0].URL, tt.url)
				}
				
				if feeds[0].Title != tt.title {
					t.Errorf("Title mismatch: got %q, want %q", feeds[0].Title, tt.title)
				}
			}
		})
	}
}

// TestUpdateFeedEdgeCases tests update with edge cases
func TestUpdateFeedEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		initialURL  string
		initialTitle string
		updateURL   string
		updateTitle string
		updateID    int
		wantErr     bool
	}{
		{
			name:        "Update non-existent feed",
			initialURL:  "https://example.com/exists.xml",
			initialTitle: "Exists",
			updateURL:   "https://example.com/new.xml",
			updateTitle: "Updated",
			updateID:    99999, // Non-existent ID
			wantErr:     false, // SQLite silently succeeds if no rows affected
		},
		{
			name:        "Update with negative ID",
			initialURL:  "https://example.com/exists.xml",
			initialTitle: "Exists",
			updateURL:   "https://example.com/new.xml",
			updateTitle: "Updated",
			updateID:    -1,
			wantErr:     false,
		},
		{
			name:        "Update with zero ID",
			initialURL:  "https://example.com/exists.xml",
			initialTitle: "Exists",
			updateURL:   "https://example.com/new.xml",
			updateTitle: "Updated",
			updateID:    0,
			wantErr:     false,
		},
		{
			name:        "Update to empty strings",
			initialURL:  "https://example.com/exists.xml",
			initialTitle: "Exists",
			updateURL:   "",
			updateTitle: "",
			updateID:    1,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			
			// Add initial feed
			if err := db.AddFeed(tt.initialURL, tt.initialTitle); err != nil {
				t.Fatalf("failed to add initial feed: %v", err)
			}

			// Update feed
			feed := models.Feed{
				ID:    tt.updateID,
				URL:   tt.updateURL,
				Title: tt.updateTitle,
			}
			
			err := db.UpdateFeed(feed)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateFeed() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestDeleteFeedEdgeCases tests delete with edge cases
func TestDeleteFeedEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "Delete non-existent feed",
			id:      99999,
			wantErr: false, // SQLite silently succeeds
		},
		{
			name:    "Delete with negative ID",
			id:      -1,
			wantErr: false,
		},
		{
			name:    "Delete with zero ID",
			id:      0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			
			err := db.DeleteFeed(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteFeed() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestConcurrentOperations tests basic concurrent safety
func TestConcurrentOperations(t *testing.T) {
	db := setupTestDB(t)
	
	// Add initial feeds
	for i := 0; i < 10; i++ {
		url := fmt.Sprintf("https://example.com/feed%d.xml", i)
		title := fmt.Sprintf("Feed %d", i)
		if err := db.AddFeed(url, title); err != nil {
			t.Fatalf("failed to add feed %d: %v", i, err)
		}
	}

	// Concurrent reads
	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func() {
			_, err := db.GetFeeds()
			if err != nil {
				t.Errorf("concurrent GetFeeds() error: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}
}

// TestFeedTitleWithEmptyURL verifies behavior when adding a feed with an empty URL.
func TestFeedTitleWithEmptyURL(t *testing.T) {
	db := setupTestDB(t)

	// We expect this to succeed since URL column is TEXT and empty string is valid.
	// The UNIQUE constraint only prevents duplicates.
	err := db.AddFeed("", "No URL Feed")
	if err != nil {
		t.Fatalf("AddFeed with empty URL returned unexpected error: %v", err)
	}

	feeds, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("GetFeeds after inserting empty URL feed returned error: %v", err)
	}

	found := false
	for _, f := range feeds {
		if f.URL == "" && f.Title == "No URL Feed" {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected to find feed with empty URL and title %q in database", "No URL Feed")
	}
}
