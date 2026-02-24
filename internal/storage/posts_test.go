package storage

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/pixel-87/warss/internal/models"
)

// TestAddPosts tests the AddPosts method with various scenarios
func TestAddPosts(t *testing.T) {
	tests := []struct {
		name    string
		feedID  int
		posts   []models.Post
		wantErr bool
	}{
		{
			name:   "Add single valid post",
			feedID: 1,
			posts: []models.Post{
				{
					Title:       "Test Post",
					Link:        "http://example.com/post1",
					Content:     "Test content",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name:   "Add multiple posts",
			feedID: 1,
			posts: []models.Post{
				{
					Title:       "Post 1",
					Link:        "http://example.com/post2",
					Content:     "Content 1",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					Title:       "Post 2",
					Link:        "http://example.com/post3",
					Content:     "Content 2",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name:    "Add empty post slice",
			feedID:  1,
			posts:   []models.Post{},
			wantErr: false,
		},
		{
			name:    "Add nil post slice",
			feedID:  1,
			posts:   nil,
			wantErr: false,
		},
		{
			name:   "Add post with empty content",
			feedID: 1,
			posts: []models.Post{
				{
					Title:       "No Content",
					Link:        "http://example.com/post4",
					Content:     "",
					PublishedAt: time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name:   "Add post with zero times",
			feedID: 1,
			posts: []models.Post{
				{
					Title:       "Zero Times",
					Link:        "http://example.com/post5",
					Content:     "Content",
					PublishedAt: time.Time{},
					UpdatedAt:   time.Time{},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)

			// Add a feed first
			if err := db.AddFeed("http://example.com/feed.xml", "Test Feed"); err != nil {
				t.Fatalf("failed to add test feed: %v", err)
			}

			err := db.AddPosts(tt.feedID, tt.posts)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddPosts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestAddPostsConflict tests the ON CONFLICT DO NOTHING behavior
func TestAddPostsConflict(t *testing.T) {
	db := setupTestDB(t)

	// Add a feed
	if err := db.AddFeed("http://example.com/feed.xml", "Test Feed"); err != nil {
		t.Fatalf("failed to add test feed: %v", err)
	}

	now := time.Now()
	post := models.Post{
		Title:       "Duplicate Post",
		Link:        "http://example.com/duplicate",
		Content:     "Original Content",
		PublishedAt: now,
		UpdatedAt:   now,
	}

	// Add post first time
	if err := db.AddPosts(1, []models.Post{post}); err != nil {
		t.Fatalf("first AddPosts() failed: %v", err)
	}

	// Try to add same link again (should be silently ignored)
	post.Title = "Updated Title"
	post.Content = "Updated Content"
	if err := db.AddPosts(1, []models.Post{post}); err != nil {
		t.Errorf("AddPosts() with duplicate link should not error: %v", err)
	}

	// Verify only one post exists
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM posts WHERE link = ?", post.Link).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count posts: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 post with duplicate link, got %d", count)
	}

	// Verify original content is preserved
	var storedTitle string
	err = db.conn.QueryRow("SELECT title FROM posts WHERE link = ?", post.Link).Scan(&storedTitle)
	if err != nil {
		t.Fatalf("failed to get stored post: %v", err)
	}

	if storedTitle != "Duplicate Post" {
		t.Errorf("expected original title %q, got %q", "Duplicate Post", storedTitle)
	}
}

// TestAddPostsInvalidFeedID tests adding posts with invalid feed IDs
func TestAddPostsInvalidFeedID(t *testing.T) {
	tests := []struct {
		name    string
		feedID  int
		wantErr bool
	}{
		{
			name:    "Non-existent feed ID",
			feedID:  99999,
			wantErr: true, // Foreign key constraint
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
			db := setupTestDB(t)

			post := models.Post{
				Title:       "Test",
				Link:        "http://example.com/test",
				Content:     "Content",
				PublishedAt: time.Now(),
				UpdatedAt:   time.Now(),
			}

			err := db.AddPosts(tt.feedID, []models.Post{post})
			if (err != nil) != tt.wantErr {
				t.Errorf("AddPosts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestAddPostsEdgeCases tests edge cases with unusual post data
func TestAddPostsEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		post    models.Post
		wantErr bool
	}{
		{
			name: "Very long title",
			post: models.Post{
				Title:       strings.Repeat("A", 10000),
				Link:        "http://example.com/long-title",
				Content:     "Content",
				PublishedAt: time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Very long content",
			post: models.Post{
				Title:       "Title",
				Link:        "http://example.com/long-content",
				Content:     strings.Repeat("Lorem ipsum ", 100000),
				PublishedAt: time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Unicode in all fields",
			post: models.Post{
				Title:       "日本語のタイトル 🚀",
				Link:        "http://example.com/日本語",
				Content:     "محتوى عربي with emojis 🎉🎊",
				PublishedAt: time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "HTML in content",
			post: models.Post{
				Title:       "HTML Post",
				Link:        "http://example.com/html",
				Content:     "<script>alert('xss')</script><p>Paragraph</p>",
				PublishedAt: time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Special characters in link",
			post: models.Post{
				Title:       "Query Params",
				Link:        "http://example.com/path?param=value&other=test#anchor",
				Content:     "Content",
				PublishedAt: time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Newlines and tabs in fields",
			post: models.Post{
				Title:       "Multi-line\nTitle\tWith\tTabs",
				Link:        "http://example.com/multiline",
				Content:     "Line 1\nLine 2\tTabbed\rCarriage Return",
				PublishedAt: time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "SQL injection attempt in title",
			post: models.Post{
				Title:       "'; DROP TABLE posts; --",
				Link:        "http://example.com/sql-inject",
				Content:     "Content",
				PublishedAt: time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Empty title (should still work)",
			post: models.Post{
				Title:       "",
				Link:        "http://example.com/empty-title",
				Content:     "Content",
				PublishedAt: time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Whitespace only in fields",
			post: models.Post{
				Title:       "   \t\n   ",
				Link:        "http://example.com/whitespace",
				Content:     "   \t\n   ",
				PublishedAt: time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Very long URL",
			post: models.Post{
				Title:       "Long URL",
				Link:        "http://example.com/" + strings.Repeat("path/", 500),
				Content:     "Content",
				PublishedAt: time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)

			// Add a feed
			if err := db.AddFeed("http://example.com/feed.xml", "Test Feed"); err != nil {
				t.Fatalf("failed to add test feed: %v", err)
			}

			err := db.AddPosts(1, []models.Post{tt.post})
			if (err != nil) != tt.wantErr {
				t.Errorf("AddPosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify post was stored correctly if no error expected
			if !tt.wantErr {
				var count int
				err := db.conn.QueryRow("SELECT COUNT(*) FROM posts WHERE link = ?", tt.post.Link).Scan(&count)
				if err != nil {
					t.Fatalf("failed to verify post: %v", err)
				}
				if count != 1 {
					t.Errorf("expected post to be stored, got count %d", count)
				}
			}
		})
	}
}

// TestGetPost tests the GetPost method
func TestGetPost(t *testing.T) {
	db := setupTestDB(t)

	// Add a feed
	if err := db.AddFeed("http://example.com/feed.xml", "Test Feed"); err != nil {
		t.Fatalf("failed to add test feed: %v", err)
	}

	now := time.Now().Truncate(time.Second)
	testPost := models.Post{
		Title:       "Test Post",
		Link:        "http://example.com/post",
		Content:     "Test Content",
		PublishedAt: now,
		UpdatedAt:   now,
	}

	// Add post
	if err := db.AddPosts(1, []models.Post{testPost}); err != nil {
		t.Fatalf("failed to add test post: %v", err)
	}

	// Get the post ID
	var postID int
	err := db.conn.QueryRow("SELECT id FROM posts WHERE link = ?", testPost.Link).Scan(&postID)
	if err != nil {
		t.Fatalf("failed to get post ID: %v", err)
	}

	// Test getting the post
	ctx := context.Background()
	retrieved, err := db.GetPost(ctx, postID)
	if err != nil {
		t.Fatalf("GetPost() error = %v", err)
	}

	// Verify fields
	if retrieved.Title != testPost.Title {
		t.Errorf("Title = %q, want %q", retrieved.Title, testPost.Title)
	}
	if retrieved.Link != testPost.Link {
		t.Errorf("Link = %q, want %q", retrieved.Link, testPost.Link)
	}
	if retrieved.Content != testPost.Content {
		t.Errorf("Content = %q, want %q", retrieved.Content, testPost.Content)
	}
	if retrieved.FeedID != 1 {
		t.Errorf("FeedID = %d, want 1", retrieved.FeedID)
	}
}

// TestGetPostNotFound tests GetPost with non-existent post IDs
func TestGetPostNotFound(t *testing.T) {
	tests := []struct {
		name    string
		postID  int
		wantErr bool
	}{
		{
			name:    "Non-existent post ID",
			postID:  99999,
			wantErr: true,
		},
		{
			name:    "Negative post ID",
			postID:  -1,
			wantErr: true,
		},
		{
			name:    "Zero post ID",
			postID:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)

			ctx := context.Background()
			_, err := db.GetPost(ctx, tt.postID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPost() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify it's a sql.ErrNoRows error if expected
			if tt.wantErr && err != nil && err != sql.ErrNoRows && !strings.Contains(err.Error(), "no rows") {
				t.Logf("Expected sql.ErrNoRows or similar, got: %v", err)
			}
		})
	}
}

// TestGetPostWithContext tests GetPost with various context scenarios
func TestGetPostWithContext(t *testing.T) {
	db := setupTestDB(t)

	// Add a feed and post
	if err := db.AddFeed("http://example.com/feed.xml", "Test Feed"); err != nil {
		t.Fatalf("failed to add test feed: %v", err)
	}

	testPost := models.Post{
		Title:       "Context Test Post",
		Link:        "http://example.com/context",
		Content:     "Content",
		PublishedAt: time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.AddPosts(1, []models.Post{testPost}); err != nil {
		t.Fatalf("failed to add test post: %v", err)
	}

	var postID int
	err := db.conn.QueryRow("SELECT id FROM posts WHERE link = ?", testPost.Link).Scan(&postID)
	if err != nil {
		t.Fatalf("failed to get post ID: %v", err)
	}

	t.Run("Normal context", func(t *testing.T) {
		ctx := context.Background()
		_, err := db.GetPost(ctx, postID)
		if err != nil {
			t.Errorf("GetPost() with normal context error = %v", err)
		}
	})

	t.Run("Context with timeout (not triggered)", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := db.GetPost(ctx, postID)
		if err != nil {
			t.Errorf("GetPost() with timeout context error = %v", err)
		}
	})

	t.Run("Cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := db.GetPost(ctx, postID)
		// Should get context cancelled error
		if err == nil {
			t.Error("GetPost() with cancelled context should return error")
		}
	})
}

// TestAddPostsMultipleBatches tests adding posts in multiple batches
func TestAddPostsMultipleBatches(t *testing.T) {
	db := setupTestDB(t)

	// Add a feed
	if err := db.AddFeed("http://example.com/feed.xml", "Test Feed"); err != nil {
		t.Fatalf("failed to add test feed: %v", err)
	}

	// Add first batch
	batch1 := []models.Post{
		{
			Title:       "Post 1",
			Link:        "http://example.com/1",
			Content:     "Content 1",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Title:       "Post 2",
			Link:        "http://example.com/2",
			Content:     "Content 2",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	if err := db.AddPosts(1, batch1); err != nil {
		t.Fatalf("failed to add batch 1: %v", err)
	}

	// Add second batch
	batch2 := []models.Post{
		{
			Title:       "Post 3",
			Link:        "http://example.com/3",
			Content:     "Content 3",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Title:       "Post 4",
			Link:        "http://example.com/4",
			Content:     "Content 4",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	if err := db.AddPosts(1, batch2); err != nil {
		t.Fatalf("failed to add batch 2: %v", err)
	}

	// Verify total count
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM posts WHERE feed_id = 1").Scan(&count)
	if err != nil {
		t.Fatalf("failed to count posts: %v", err)
	}

	if count != 4 {
		t.Errorf("expected 4 posts, got %d", count)
	}
}

// TestAddPostsLargeVolume tests adding a large number of posts
func TestAddPostsLargeVolume(t *testing.T) {
	db := setupTestDB(t)

	// Add a feed
	if err := db.AddFeed("http://example.com/feed.xml", "Test Feed"); err != nil {
		t.Fatalf("failed to add test feed: %v", err)
	}

	// Create 100 posts
	posts := make([]models.Post, 100)
	now := time.Now()
	for i := range posts {
		posts[i] = models.Post{
			Title:       "Post " + strings.Repeat("x", i),
			Link:        "http://example.com/post" + strings.Repeat("x", i),
			Content:     "Content " + strings.Repeat("y", i),
			PublishedAt: now.Add(time.Duration(i) * time.Minute),
			UpdatedAt:   now.Add(time.Duration(i) * time.Minute),
		}
	}

	if err := db.AddPosts(1, posts); err != nil {
		t.Fatalf("AddPosts() with large volume error = %v", err)
	}

	// Verify count
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM posts WHERE feed_id = 1").Scan(&count)
	if err != nil {
		t.Fatalf("failed to count posts: %v", err)
	}

	if count != 100 {
		t.Errorf("expected 100 posts, got %d", count)
	}
}

// TestGetPostDataIntegrity verifies data integrity after round-trip
func TestGetPostDataIntegrity(t *testing.T) {
	db := setupTestDB(t)

	// Add a feed
	if err := db.AddFeed("http://example.com/feed.xml", "Test Feed"); err != nil {
		t.Fatalf("failed to add test feed: %v", err)
	}

	// Create a post with specific data
	now := time.Date(2024, 1, 15, 12, 30, 45, 0, time.UTC)
	originalPost := models.Post{
		Title:       "Data Integrity Test 🔒",
		Link:        "http://example.com/integrity",
		Content:     "Line 1\nLine 2\tTabbed\r\nSpecial: <>&\"'",
		PublishedAt: now,
		UpdatedAt:   now.Add(1 * time.Hour),
	}

	// Add the post
	if err := db.AddPosts(1, []models.Post{originalPost}); err != nil {
		t.Fatalf("failed to add post: %v", err)
	}

	// Get the post ID
	var postID int
	err := db.conn.QueryRow("SELECT id FROM posts WHERE link = ?", originalPost.Link).Scan(&postID)
	if err != nil {
		t.Fatalf("failed to get post ID: %v", err)
	}

	// Retrieve the post
	ctx := context.Background()
	retrieved, err := db.GetPost(ctx, postID)
	if err != nil {
		t.Fatalf("GetPost() error = %v", err)
	}

	// Verify exact match
	if retrieved.Title != originalPost.Title {
		t.Errorf("Title mismatch:\ngot:  %q\nwant: %q", retrieved.Title, originalPost.Title)
	}
	if retrieved.Link != originalPost.Link {
		t.Errorf("Link mismatch:\ngot:  %q\nwant: %q", retrieved.Link, originalPost.Link)
	}
	if retrieved.Content != originalPost.Content {
		t.Errorf("Content mismatch:\ngot:  %q\nwant: %q", retrieved.Content, originalPost.Content)
	}
	if retrieved.FeedID != 1 {
		t.Errorf("FeedID = %d, want 1", retrieved.FeedID)
	}

	// Note: Time comparison might have precision issues depending on SQLite driver
	// Truncate to second for comparison
	if !retrieved.PublishedAt.Truncate(time.Second).Equal(originalPost.PublishedAt.Truncate(time.Second)) {
		t.Errorf("PublishedAt mismatch:\ngot:  %v\nwant: %v", retrieved.PublishedAt, originalPost.PublishedAt)
	}
	if !retrieved.UpdatedAt.Truncate(time.Second).Equal(originalPost.UpdatedAt.Truncate(time.Second)) {
		t.Errorf("UpdatedAt mismatch:\ngot:  %v\nwant: %v", retrieved.UpdatedAt, originalPost.UpdatedAt)
	}
}

// TestAddPostsPartialFailure tests behavior when one post in a batch fails
func TestAddPostsPartialFailure(t *testing.T) {
	db := setupTestDB(t)

	// Add a feed
	if err := db.AddFeed("http://example.com/feed.xml", "Test Feed"); err != nil {
		t.Fatalf("failed to add test feed: %v", err)
	}

	// Add a post first
	firstPost := models.Post{
		Title:       "First",
		Link:        "http://example.com/duplicate",
		Content:     "Content",
		PublishedAt: time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.AddPosts(1, []models.Post{firstPost}); err != nil {
		t.Fatalf("failed to add first post: %v", err)
	}

	// Try to add a batch where one has duplicate link (should be ignored via ON CONFLICT)
	batch := []models.Post{
		{
			Title:       "New Post 1",
			Link:        "http://example.com/new1",
			Content:     "Content",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Title:       "Duplicate Link",
			Link:        "http://example.com/duplicate", // Duplicate!
			Content:     "Different Content",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Title:       "New Post 2",
			Link:        "http://example.com/new2",
			Content:     "Content",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Should succeed (duplicate is ignored)
	if err := db.AddPosts(1, batch); err != nil {
		t.Fatalf("AddPosts() with duplicate should not error: %v", err)
	}

	// Verify we have 3 posts total (1 original + 2 new, 1 duplicate ignored)
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM posts WHERE feed_id = 1").Scan(&count)
	if err != nil {
		t.Fatalf("failed to count posts: %v", err)
	}

	if count != 3 {
		t.Errorf("expected 3 posts (duplicate ignored), got %d", count)
	}
}

// TestForeignKeyConstraint verifies cascade delete works
func TestForeignKeyConstraint(t *testing.T) {
	db := setupTestDB(t)

	// Add a feed
	if err := db.AddFeed("http://example.com/feed.xml", "Test Feed"); err != nil {
		t.Fatalf("failed to add test feed: %v", err)
	}

	// Add posts to the feed
	posts := []models.Post{
		{
			Title:       "Post 1",
			Link:        "http://example.com/1",
			Content:     "Content 1",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Title:       "Post 2",
			Link:        "http://example.com/2",
			Content:     "Content 2",
			PublishedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	if err := db.AddPosts(1, posts); err != nil {
		t.Fatalf("failed to add posts: %v", err)
	}

	// Verify posts exist
	var countBefore int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM posts WHERE feed_id = 1").Scan(&countBefore)
	if err != nil {
		t.Fatalf("failed to count posts before delete: %v", err)
	}

	if countBefore != 2 {
		t.Errorf("expected 2 posts before delete, got %d", countBefore)
	}

	// Delete the feed (should cascade delete posts)
	if err := db.DeleteFeed(1); err != nil {
		t.Fatalf("failed to delete feed: %v", err)
	}

	// Verify posts are gone
	var countAfter int
	err = db.conn.QueryRow("SELECT COUNT(*) FROM posts WHERE feed_id = 1").Scan(&countAfter)
	if err != nil {
		t.Fatalf("failed to count posts after delete: %v", err)
	}

	if countAfter != 0 {
		t.Errorf("expected 0 posts after cascade delete, got %d", countAfter)
	}
}
