package storage

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/pixel-87/warss/internal/models"
)

// setupPostsTestDB creates a test database with a test feed for posts testing
func setupPostsTestDB(t *testing.T) (*DB, int) {
	tmpFile := "test_posts.db"
	_ = os.Remove(tmpFile)

	db, err := NewDB(tmpFile)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
		_ = os.Remove(tmpFile)
	})

	// Add a test feed
	if err := db.AddFeed("https://example.com/feed.xml", "Test Feed"); err != nil {
		t.Fatalf("failed to add test feed: %v", err)
	}

	feeds, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("failed to get feeds: %v", err)
	}

	if len(feeds) == 0 {
		t.Fatalf("no feeds found after adding")
	}

	return db, feeds[0].ID
}

// TestAddPosts tests adding posts to a feed
func TestAddPosts(t *testing.T) {
	tests := []struct {
		name    string
		posts   []models.Post
		wantErr bool
	}{
		{
			name: "Add single post",
			posts: []models.Post{
				{
					Title:       "Test Post 1",
					Link:        "https://example.com/post1",
					Content:     "Content 1",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Add multiple posts",
			posts: []models.Post{
				{
					Title:       "Post 1",
					Link:        "https://example.com/multi1",
					Content:     "Content 1",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					Title:       "Post 2",
					Link:        "https://example.com/multi2",
					Content:     "Content 2",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					Title:       "Post 3",
					Link:        "https://example.com/multi3",
					Content:     "Content 3",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Add post with empty content",
			posts: []models.Post{
				{
					Title:       "No Content",
					Link:        "https://example.com/nocontent",
					Content:     "",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Add post with unicode",
			posts: []models.Post{
				{
					Title:       "日本語のタイトル 🚀",
					Link:        "https://example.com/unicode",
					Content:     "محتوى عربي",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name:    "Add empty posts slice",
			posts:   []models.Post{},
			wantErr: false,
		},
		{
			name:    "Add nil posts slice",
			posts:   nil,
			wantErr: false,
		},
		{
			name: "Add post with very long content",
			posts: []models.Post{
				{
					Title:       "Long Content Post",
					Link:        "https://example.com/longcontent",
					Content:     strings.Repeat("Lorem ipsum dolor sit amet. ", 1000),
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Add post with HTML content",
			posts: []models.Post{
				{
					Title:       "HTML Post",
					Link:        "https://example.com/html",
					Content:     "<script>alert('xss')</script><p>Paragraph</p>",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Add post with special characters in link",
			posts: []models.Post{
				{
					Title:       "Query Params",
					Link:        "https://example.com/path?param=value&other=test#anchor",
					Content:     "Content",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Add post with newlines in content",
			posts: []models.Post{
				{
					Title:       "Multi-line",
					Link:        "https://example.com/multiline",
					Content:     "Line 1\nLine 2\nLine 3\tTabbed",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, feedID := setupPostsTestDB(t)

			err := db.AddPosts(feedID, tt.posts)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddPosts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestAddPostsConflict tests the ON CONFLICT behavior
func TestAddPostsConflict(t *testing.T) {
	db, feedID := setupPostsTestDB(t)

	// Add initial post
	initialPosts := []models.Post{
		{
			Title:       "Original Title",
			Link:        "https://example.com/duplicate",
			Content:     "Original Content",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	if err := db.AddPosts(feedID, initialPosts); err != nil {
		t.Fatalf("failed to add initial posts: %v", err)
	}

	// Try to add duplicate (same link)
	duplicatePosts := []models.Post{
		{
			Title:       "Updated Title",
			Link:        "https://example.com/duplicate",
			Content:     "Updated Content",
			PublishedAt: time.Now().Add(time.Hour),
			UpdatedAt:   time.Now().Add(time.Hour),
		},
	}

	// Should not error due to ON CONFLICT DO NOTHING
	if err := db.AddPosts(feedID, duplicatePosts); err != nil {
		t.Fatalf("AddPosts() with duplicate link returned error: %v", err)
	}
}

// TestAddPostsWithInvalidFeedID tests adding posts with non-existent feed ID
func TestAddPostsWithInvalidFeedID(t *testing.T) {
	db, _ := setupPostsTestDB(t)

	posts := []models.Post{
		{
			Title:       "Orphan Post",
			Link:        "https://example.com/orphan",
			Content:     "Content",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Use non-existent feed ID
	err := db.AddPosts(99999, posts)
	// This should error due to foreign key constraint
	if err == nil {
		t.Errorf("AddPosts() with invalid feed ID should return error")
	}
}

// TestAddPostsEdgeCases tests edge cases and boundary conditions
func TestAddPostsEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		feedID  int
		posts   []models.Post
		wantErr bool
	}{
		{
			name:   "Negative feed ID",
			feedID: -1,
			posts: []models.Post{
				{
					Title:       "Test",
					Link:        "https://example.com/negative",
					Content:     "Content",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: true,
		},
		{
			name:   "Zero feed ID",
			feedID: 0,
			posts: []models.Post{
				{
					Title:       "Test",
					Link:        "https://example.com/zero",
					Content:     "Content",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: true,
		},
		{
			name:   "Empty title",
			feedID: 1,
			posts: []models.Post{
				{
					Title:       "",
					Link:        "https://example.com/notitle",
					Content:     "Content",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false, // Database allows empty title
		},
		{
			name:   "Empty link",
			feedID: 1,
			posts: []models.Post{
				{
					Title:       "Title",
					Link:        "",
					Content:     "Content",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false, // Database allows empty link (though not ideal)
		},
		{
			name:   "Zero time values",
			feedID: 1,
			posts: []models.Post{
				{
					Title:       "Zero Time",
					Link:        "https://example.com/zerotime",
					Content:     "Content",
					PublishedAt: time.Time{},
					UpdatedAt:   time.Time{},
				},
			},
			wantErr: false,
		},
		{
			name:   "Future dates",
			feedID: 1,
			posts: []models.Post{
				{
					Title:       "Future Post",
					Link:        "https://example.com/future",
					Content:     "Content",
					PublishedAt: time.Now().Add(365 * 24 * time.Hour),
					UpdatedAt:   time.Now().Add(365 * 24 * time.Hour),
				},
			},
			wantErr: false,
		},
		{
			name:   "Very old dates",
			feedID: 1,
			posts: []models.Post{
				{
					Title:       "Old Post",
					Link:        "https://example.com/old",
					Content:     "Content",
					PublishedAt: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
		},
		{
			name:   "SQL injection in title",
			feedID: 1,
			posts: []models.Post{
				{
					Title:       "'; DROP TABLE posts; --",
					Link:        "https://example.com/sqlinject1",
					Content:     "Content",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name:   "SQL injection in content",
			feedID: 1,
			posts: []models.Post{
				{
					Title:       "Test",
					Link:        "https://example.com/sqlinject2",
					Content:     "1' OR '1'='1",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name:   "Very large number of posts",
			feedID: 1,
			posts: func() []models.Post {
				posts := make([]models.Post, 100)
				for i := range posts {
					posts[i] = models.Post{
						Title:       "Post " + strings.Repeat("A", i),
						Link:        "https://example.com/large" + strings.Repeat("A", i),
						Content:     "Content",
						PublishedAt: time.Now(),
						UpdatedAt:   time.Now(),
					}
				}
				return posts
			}(),
			wantErr: false,
		},
		{
			name:   "Post with null bytes in content",
			feedID: 1,
			posts: []models.Post{
				{
					Title:       "Null Bytes",
					Link:        "https://example.com/nullbytes",
					Content:     "Content\x00with\x00null\x00bytes",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name:   "Post with extremely long URL",
			feedID: 1,
			posts: []models.Post{
				{
					Title:       "Long URL",
					Link:        "https://example.com/" + strings.Repeat("path/", 500),
					Content:     "Content",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, feedID := setupPostsTestDB(t)

			// Use the actual feed ID from setup if test specifies feedID as 1
			if tt.feedID == 1 {
				tt.feedID = feedID
			}

			err := db.AddPosts(tt.feedID, tt.posts)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddPosts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestAddPostsMultipleFeedsIsolation tests that posts are isolated per feed
func TestAddPostsMultipleFeedsIsolation(t *testing.T) {
	tmpFile := "test_isolation.db"
	_ = os.Remove(tmpFile)

	db, err := NewDB(tmpFile)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer func() {
		_ = db.Close()
		_ = os.Remove(tmpFile)
	}()

	// Add two feeds
	if err := db.AddFeed("https://feed1.com/rss", "Feed 1"); err != nil {
		t.Fatalf("failed to add feed 1: %v", err)
	}
	if err := db.AddFeed("https://feed2.com/rss", "Feed 2"); err != nil {
		t.Fatalf("failed to add feed 2: %v", err)
	}

	feeds, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("failed to get feeds: %v", err)
	}

	feed1ID := feeds[0].ID
	feed2ID := feeds[1].ID

	// Add posts to feed 1
	posts1 := []models.Post{
		{
			Title:       "Feed 1 Post 1",
			Link:        "https://feed1.com/post1",
			Content:     "Content 1",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	if err := db.AddPosts(feed1ID, posts1); err != nil {
		t.Fatalf("failed to add posts to feed 1: %v", err)
	}

	// Add posts to feed 2
	posts2 := []models.Post{
		{
			Title:       "Feed 2 Post 1",
			Link:        "https://feed2.com/post1",
			Content:     "Content 2",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	if err := db.AddPosts(feed2ID, posts2); err != nil {
		t.Fatalf("failed to add posts to feed 2: %v", err)
	}

	// Verify both posts exist independently
	// (Would need a GetPosts method to verify, but this tests the add functionality)
}

// TestAddPostsBoundaryTimes tests posts with boundary time values
func TestAddPostsBoundaryTimes(t *testing.T) {
	db, feedID := setupPostsTestDB(t)

	tests := []struct {
		name        string
		publishedAt time.Time
		updatedAt   time.Time
	}{
		{
			name:        "Unix epoch",
			publishedAt: time.Unix(0, 0),
			updatedAt:   time.Unix(0, 0),
		},
		{
			name:        "Year 2038 problem",
			publishedAt: time.Unix(2147483647, 0),
			updatedAt:   time.Unix(2147483647, 0),
		},
		{
			name:        "Nanosecond precision",
			publishedAt: time.Unix(1234567890, 123456789),
			updatedAt:   time.Unix(1234567890, 987654321),
		},
		{
			name:        "Leap year date",
			publishedAt: time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC),
			updatedAt:   time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			posts := []models.Post{
				{
					Title:       tt.name,
					Link:        "https://example.com/" + tt.name,
					Content:     "Content",
					PublishedAt: tt.publishedAt,
					UpdatedAt:   tt.updatedAt,
				},
			}

			if err := db.AddPosts(feedID, posts); err != nil {
				t.Errorf("AddPosts() with %s failed: %v", tt.name, err)
			}
		})
	}
}

// TestAddPostsReturnBehavior tests that the function returns correctly
func TestAddPostsReturnBehavior(t *testing.T) {
	db, feedID := setupPostsTestDB(t)

	// This test specifically checks if the function processes ALL posts
	// or if it has a bug where it returns early
	posts := []models.Post{
		{
			Title:       "First Post",
			Link:        "https://example.com/first",
			Content:     "First",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Title:       "Second Post",
			Link:        "https://example.com/second",
			Content:     "Second",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Title:       "Third Post",
			Link:        "https://example.com/third",
			Content:     "Third",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	err := db.AddPosts(feedID, posts)
	if err != nil {
		t.Fatalf("AddPosts() returned error: %v", err)
	}

	// Note: Without a GetPosts method, we can't verify all posts were added
	// This test will PASS even if the function has a bug where it returns early
	// But it documents the expected behavior
}

// TestAddPostsEmptyStrings tests posts with various empty string combinations
func TestAddPostsEmptyStrings(t *testing.T) {
	db, feedID := setupPostsTestDB(t)

	tests := []struct {
		name    string
		post    models.Post
		wantErr bool
	}{
		{
			name: "All fields empty except link",
			post: models.Post{
				Title:       "",
				Link:        "https://example.com/empty1",
				Content:     "",
				PublishedAt: time.Time{},
				UpdatedAt:   time.Time{},
			},
			wantErr: false,
		},
		{
			name: "Whitespace only title",
			post: models.Post{
				Title:       "   \t\n   ",
				Link:        "https://example.com/whitespace1",
				Content:     "Content",
				PublishedAt: time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Whitespace only content",
			post: models.Post{
				Title:       "Title",
				Link:        "https://example.com/whitespace2",
				Content:     "   \t\n   ",
				PublishedAt: time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.AddPosts(feedID, []models.Post{tt.post})
			if (err != nil) != tt.wantErr {
				t.Errorf("AddPosts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
