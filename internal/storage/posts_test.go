package storage

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/pixel-87/warss/internal/models"
)

// TestAddPosts tests the AddPosts method
func TestAddPosts(t *testing.T) {
	db := setupTestDB(t)

	// Create a feed first
	if err := db.AddFeed("https://example.com/feed.xml", "Test Feed"); err != nil {
		t.Fatalf("failed to add test feed: %v", err)
	}

	feeds, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("failed to get feeds: %v", err)
	}
	feedID := feeds[0].ID

	now := time.Now()

	tests := []struct {
		name    string
		feedID  int
		posts   []models.Post
		wantErr bool
	}{
		{
			name:   "Add single post",
			feedID: feedID,
			posts: []models.Post{
				{
					Title:       "Test Post",
					Link:        "https://example.com/post1",
					Content:     "Test Content",
					PublishedAt: now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
		},
		{
			name:   "Add multiple posts",
			feedID: feedID,
			posts: []models.Post{
				{
					Title:       "Post 2",
					Link:        "https://example.com/post2",
					Content:     "Content 2",
					PublishedAt: now,
					UpdatedAt:   now,
				},
				{
					Title:       "Post 3",
					Link:        "https://example.com/post3",
					Content:     "Content 3",
					PublishedAt: now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
		},
		{
			name:   "Add duplicate post (same link)",
			feedID: feedID,
			posts: []models.Post{
				{
					Title:       "Duplicate",
					Link:        "https://example.com/post1", // Same as first test
					Content:     "Should be ignored",
					PublishedAt: now,
					UpdatedAt:   now,
				},
			},
			wantErr: false, // ON CONFLICT DO NOTHING
		},
		{
			name:    "Add empty posts slice",
			feedID:  feedID,
			posts:   []models.Post{},
			wantErr: false,
		},
		{
			name:    "Add nil posts slice",
			feedID:  feedID,
			posts:   nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.AddPosts(tt.feedID, tt.posts)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddPosts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestAddPostsEdgeCases tests edge cases and unusual inputs
func TestAddPostsEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		posts   []models.Post
		wantErr bool
	}{
		{
			name: "Very long title",
			posts: []models.Post{
				{
					Title:       strings.Repeat("A", 10000),
					Link:        "https://example.com/long-title",
					Content:     "Content",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Very long content",
			posts: []models.Post{
				{
					Title:       "Post",
					Link:        "https://example.com/long-content",
					Content:     strings.Repeat("Lorem ipsum ", 10000),
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Unicode in all fields",
			posts: []models.Post{
				{
					Title:       "日本語のタイトル 🚀",
					Link:        "https://example.com/日本語",
					Content:     "محتوى عربي with 中文 and emoji 🎉",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "HTML in content",
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
			name: "Empty title",
			posts: []models.Post{
				{
					Title:       "",
					Link:        "https://example.com/empty-title",
					Content:     "Content",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Empty content",
			posts: []models.Post{
				{
					Title:       "Title",
					Link:        "https://example.com/empty-content",
					Content:     "",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Whitespace only fields",
			posts: []models.Post{
				{
					Title:       "   \t\n   ",
					Link:        "https://example.com/whitespace",
					Content:     "   \n\n\t   ",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Special characters in link",
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
			name: "SQL injection attempt in title",
			posts: []models.Post{
				{
					Title:       "'; DROP TABLE posts; --",
					Link:        "https://example.com/sql-injection",
					Content:     "Content",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Newlines and tabs in content",
			posts: []models.Post{
				{
					Title:       "Multi-line",
					Link:        "https://example.com/multiline",
					Content:     "Line 1\nLine 2\tTabbed\rCarriage Return",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Zero time values",
			posts: []models.Post{
				{
					Title:       "Zero Time",
					Link:        "https://example.com/zero-time",
					Content:     "Content",
					PublishedAt: time.Time{},
					UpdatedAt:   time.Time{},
				},
			},
			wantErr: false,
		},
		{
			name: "Future dates",
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
			name: "Very old dates",
			posts: []models.Post{
				{
					Title:       "Old Post",
					Link:        "https://example.com/old",
					Content:     "Content",
					PublishedAt: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)

			// Create a feed
			if err := db.AddFeed("https://example.com/feed.xml", "Test Feed"); err != nil {
				t.Fatalf("failed to add test feed: %v", err)
			}

			feeds, err := db.GetFeeds()
			if err != nil {
				t.Fatalf("failed to get feeds: %v", err)
			}
			feedID := feeds[0].ID

			err = db.AddPosts(feedID, tt.posts)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddPosts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestAddPostsInvalidFeedID tests adding posts with invalid feed IDs
func TestAddPostsInvalidFeedID(t *testing.T) {
	db := setupTestDB(t)

	now := time.Now()
	posts := []models.Post{
		{
			Title:       "Test Post",
			Link:        "https://example.com/test",
			Content:     "Content",
			PublishedAt: now,
			UpdatedAt:   now,
		},
	}

	tests := []struct {
		name    string
		feedID  int
		wantErr bool
	}{
		{
			name:    "Non-existent feed ID",
			feedID:  99999,
			wantErr: true, // FOREIGN KEY constraint
		},
		{
			name:    "Negative feed ID",
			feedID:  -1,
			wantErr: true,
		},
		{
			name:    "Zero feed ID",
			feedID:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.AddPosts(tt.feedID, posts)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddPosts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGetPost tests the GetPost method
func TestGetPost(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Create a feed
	if err := db.AddFeed("https://example.com/feed.xml", "Test Feed"); err != nil {
		t.Fatalf("failed to add test feed: %v", err)
	}

	feeds, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("failed to get feeds: %v", err)
	}
	feedID := feeds[0].ID

	// Add posts
	now := time.Now().Truncate(time.Second) // Truncate for comparison
	posts := []models.Post{
		{
			Title:       "Test Post 1",
			Link:        "https://example.com/post1",
			Content:     "Content 1",
			PublishedAt: now,
			UpdatedAt:   now,
		},
		{
			Title:       "Test Post 2",
			Link:        "https://example.com/post2",
			Content:     "Content 2",
			PublishedAt: now.Add(-24 * time.Hour),
			UpdatedAt:   now,
		},
	}

	if err := db.AddPosts(feedID, posts); err != nil {
		t.Fatalf("failed to add posts: %v", err)
	}

	// Get the post ID (assuming auto-increment starts at 1)
	postID := 1

	// Test getting the post
	post, err := db.GetPost(ctx, postID)
	if err != nil {
		t.Fatalf("GetPost() error = %v", err)
	}

	// Verify fields
	if post.Title != posts[0].Title {
		t.Errorf("Title = %q, want %q", post.Title, posts[0].Title)
	}
	if post.Link != posts[0].Link {
		t.Errorf("Link = %q, want %q", post.Link, posts[0].Link)
	}
	if post.Content != posts[0].Content {
		t.Errorf("Content = %q, want %q", post.Content, posts[0].Content)
	}
	if post.FeedID != feedID {
		t.Errorf("FeedID = %d, want %d", post.FeedID, feedID)
	}

	// Time comparison (truncate to second for SQLite)
	if !post.PublishedAt.Truncate(time.Second).Equal(posts[0].PublishedAt.Truncate(time.Second)) {
		t.Errorf("PublishedAt = %v, want %v", post.PublishedAt, posts[0].PublishedAt)
	}
	if !post.UpdatedAt.Truncate(time.Second).Equal(posts[0].UpdatedAt.Truncate(time.Second)) {
		t.Errorf("UpdatedAt = %v, want %v", post.UpdatedAt, posts[0].UpdatedAt)
	}
}

// TestGetPostNotFound tests getting a non-existent post
func TestGetPostNotFound(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	tests := []struct {
		name   string
		postID int
	}{
		{
			name:   "Non-existent post ID",
			postID: 99999,
		},
		{
			name:   "Negative post ID",
			postID: -1,
		},
		{
			name:   "Zero post ID",
			postID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := db.GetPost(ctx, tt.postID)
			if err == nil {
				t.Error("GetPost() expected error for non-existent post, got nil")
			}
			if err != sql.ErrNoRows {
				// Check that it's a "not found" type error
				if !strings.Contains(err.Error(), "no rows") {
					t.Logf("GetPost() error = %v (expected no rows error)", err)
				}
			}
		})
	}
}

// TestGetPostWithContext tests context cancellation
func TestGetPostWithContext(t *testing.T) {
	db := setupTestDB(t)

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := db.GetPost(ctx, 1)
	if err == nil {
		t.Error("GetPost() expected error with cancelled context, got nil")
	}
	if err != context.Canceled && !strings.Contains(err.Error(), "context canceled") {
		t.Logf("GetPost() with cancelled context error = %v", err)
	}
}

// TestAddPostsWithContext tests adding posts with various contexts
func TestGetPostWithTimeout(t *testing.T) {
	db := setupTestDB(t)

	// Create a feed and post
	if err := db.AddFeed("https://example.com/feed.xml", "Test Feed"); err != nil {
		t.Fatalf("failed to add test feed: %v", err)
	}

	feeds, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("failed to get feeds: %v", err)
	}
	feedID := feeds[0].ID

	posts := []models.Post{
		{
			Title:       "Test Post",
			Link:        "https://example.com/post",
			Content:     "Content",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	if err := db.AddPosts(feedID, posts); err != nil {
		t.Fatalf("failed to add posts: %v", err)
	}

	// Use a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.GetPost(ctx, 1)
	if err != nil {
		t.Errorf("GetPost() with timeout context error = %v", err)
	}
}

// TestAddPostsMultipleBatches tests adding posts in multiple batches
func TestAddPostsMultipleBatches(t *testing.T) {
	db := setupTestDB(t)

	// Create a feed
	if err := db.AddFeed("https://example.com/feed.xml", "Test Feed"); err != nil {
		t.Fatalf("failed to add test feed: %v", err)
	}

	feeds, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("failed to get feeds: %v", err)
	}
	feedID := feeds[0].ID

	// Add first batch
	batch1 := []models.Post{
		{
			Title:       "Post 1",
			Link:        "https://example.com/post1",
			Content:     "Content 1",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Title:       "Post 2",
			Link:        "https://example.com/post2",
			Content:     "Content 2",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	if err := db.AddPosts(feedID, batch1); err != nil {
		t.Fatalf("failed to add first batch: %v", err)
	}

	// Add second batch with one duplicate
	batch2 := []models.Post{
		{
			Title:       "Post 1 Duplicate",
			Link:        "https://example.com/post1", // Duplicate
			Content:     "Updated Content",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Title:       "Post 3",
			Link:        "https://example.com/post3",
			Content:     "Content 3",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	if err := db.AddPosts(feedID, batch2); err != nil {
		t.Fatalf("failed to add second batch: %v", err)
	}

	// Verify we can get the posts
	ctx := context.Background()
	post, err := db.GetPost(ctx, 1)
	if err != nil {
		t.Fatalf("GetPost() error = %v", err)
	}

	// The duplicate should have been ignored
	if post.Title != "Post 1" {
		t.Errorf("Post title = %q, want %q (duplicate should be ignored)", post.Title, "Post 1")
	}
}

// TestAddPostsLargeContent tests storing large content
func TestAddPostsLargeContent(t *testing.T) {
	db := setupTestDB(t)

	// Create a feed
	if err := db.AddFeed("https://example.com/feed.xml", "Test Feed"); err != nil {
		t.Fatalf("failed to add test feed: %v", err)
	}

	feeds, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("failed to get feeds: %v", err)
	}
	feedID := feeds[0].ID

	// Create a very large content string (1MB)
	largeContent := strings.Repeat("Lorem ipsum dolor sit amet, consectetur adipiscing elit. ", 20000)

	posts := []models.Post{
		{
			Title:       "Large Content Post",
			Link:        "https://example.com/large",
			Content:     largeContent,
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	if err := db.AddPosts(feedID, posts); err != nil {
		t.Fatalf("AddPosts() with large content error = %v", err)
	}

	// Verify it was stored correctly
	ctx := context.Background()
	post, err := db.GetPost(ctx, 1)
	if err != nil {
		t.Fatalf("GetPost() error = %v", err)
	}

	if len(post.Content) != len(largeContent) {
		t.Errorf("Content length = %d, want %d", len(post.Content), len(largeContent))
	}

	if post.Content != largeContent {
		t.Error("Large content not stored correctly")
	}
}

// TestAddPostsWithManyPosts tests adding a large number of posts
func TestAddPostsWithManyPosts(t *testing.T) {
	db := setupTestDB(t)

	// Create a feed
	if err := db.AddFeed("https://example.com/feed.xml", "Test Feed"); err != nil {
		t.Fatalf("failed to add test feed: %v", err)
	}

	feeds, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("failed to get feeds: %v", err)
	}
	feedID := feeds[0].ID

	// Create many posts
	posts := make([]models.Post, 100)
	for i := range posts {
		posts[i] = models.Post{
			Title:       "Post " + string(rune(i)),
			Link:        "https://example.com/post" + string(rune(i)),
			Content:     "Content " + string(rune(i)),
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		}
	}

	if err := db.AddPosts(feedID, posts); err != nil {
		t.Fatalf("AddPosts() with many posts error = %v", err)
	}
}
