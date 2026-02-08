package models

import (
	"strings"
	"testing"
)

// TestFeedHasUnreadPosts tests the HasUnreadPosts method
func TestFeedHasUnreadPosts(t *testing.T) {
	tests := []struct {
		name string
		feed Feed
		want bool
	}{
		{
			name: "Feed with unread posts",
			feed: Feed{
				Posts: []Post{
					{Title: "Post 1", Read: false},
					{Title: "Post 2", Read: true},
				},
			},
			want: true,
		},
		{
			name: "Feed with all posts read",
			feed: Feed{
				Posts: []Post{
					{Title: "Post 1", Read: true},
					{Title: "Post 2", Read: true},
				},
			},
			want: false,
		},
		{
			name: "Empty feed",
			feed: Feed{
				Posts: []Post{},
			},
			want: false,
		},
		{
			name: "Feed with nil posts",
			feed: Feed{
				Posts: nil,
			},
			want: false,
		},
		{
			name: "Feed with single unread post",
			feed: Feed{
				Posts: []Post{
					{Title: "Single", Read: false},
				},
			},
			want: true,
		},
		{
			name: "Feed with many posts, one unread",
			feed: Feed{
				Posts: []Post{
					{Title: "1", Read: true},
					{Title: "2", Read: true},
					{Title: "3", Read: false},
					{Title: "4", Read: true},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.feed.HasUnreadPosts()
			if got != tt.want {
				t.Errorf("HasUnreadPosts() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFeedUnreadCount tests the UnreadCount method
func TestFeedUnreadCount(t *testing.T) {
	tests := []struct {
		name string
		feed Feed
		want int
	}{
		{
			name: "Mixed read/unread posts",
			feed: Feed{
				Posts: []Post{
					{Title: "Post 1", Read: false},
					{Title: "Post 2", Read: true},
					{Title: "Post 3", Read: false},
				},
			},
			want: 2,
		},
		{
			name: "All posts unread",
			feed: Feed{
				Posts: []Post{
					{Title: "Post 1", Read: false},
					{Title: "Post 2", Read: false},
					{Title: "Post 3", Read: false},
				},
			},
			want: 3,
		},
		{
			name: "All posts read",
			feed: Feed{
				Posts: []Post{
					{Title: "Post 1", Read: true},
					{Title: "Post 2", Read: true},
				},
			},
			want: 0,
		},
		{
			name: "Empty feed",
			feed: Feed{
				Posts: []Post{},
			},
			want: 0,
		},
		{
			name: "Nil posts",
			feed: Feed{
				Posts: nil,
			},
			want: 0,
		},
		{
			name: "Large feed with many unread",
			feed: Feed{
				Posts: func() []Post {
					posts := make([]Post, 100)
					for i := 0; i < 100; i++ {
						posts[i] = Post{Read: i%2 == 0} // Half read, half unread
					}
					return posts
				}(),
			},
			want: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.feed.UnreadCount()
			if got != tt.want {
				t.Errorf("UnreadCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestPostIsValid tests the IsValid method
func TestPostIsValid(t *testing.T) {
	tests := []struct {
		name string
		post Post
		want bool
	}{
		{
			name: "Valid post with all fields",
			post: Post{
				Title:   "Valid Title",
				Content: "Valid Content",
				Link:    "http://example.com",
			},
			want: true,
		},
		{
			name: "Valid post without content",
			post: Post{
				Title: "Title",
				Link:  "http://example.com",
			},
			want: true,
		},
		{
			name: "Invalid - empty title",
			post: Post{
				Title:   "",
				Content: "Content",
				Link:    "http://example.com",
			},
			want: false,
		},
		{
			name: "Invalid - empty link",
			post: Post{
				Title:   "Title",
				Content: "Content",
				Link:    "",
			},
			want: false,
		},
		{
			name: "Invalid - both empty",
			post: Post{
				Title: "",
				Link:  "",
			},
			want: false,
		},
		{
			name: "Invalid - whitespace only title",
			post: Post{
				Title: "   \t\n  ",
				Link:  "http://example.com",
			},
			want: false,
		},
		{
			name: "Invalid - whitespace only link",
			post: Post{
				Title: "Title",
				Link:  "   \t\n  ",
			},
			want: false,
		},
		{
			name: "Valid - title and link with extra whitespace",
			post: Post{
				Title: "  Title  ",
				Link:  "  http://example.com  ",
			},
			want: true,
		},
		{
			name: "Valid - very long title",
			post: Post{
				Title:   strings.Repeat("A", 10000),
				Content: "Content",
				Link:    "http://example.com",
			},
			want: true,
		},
		{
			name: "Valid - Unicode in all fields",
			post: Post{
				Title:   "日本語のタイトル 🚀",
				Content: "محتوى عربي",
				Link:    "http://example.com/日本語",
			},
			want: true,
		},
		{
			name: "Valid - HTML in content",
			post: Post{
				Title:   "HTML Post",
				Content: "<script>alert('xss')</script><p>Paragraph</p>",
				Link:    "http://example.com",
			},
			want: true,
		},
		{
			name: "Valid - Special characters in link",
			post: Post{
				Title:   "Query Params",
				Content: "Content",
				Link:    "http://example.com/path?param=value&other=test#anchor",
			},
			want: true,
		},
		{
			name: "Valid - Newlines in title and content",
			post: Post{
				Title:   "Multi-line Title\nLine 2",
				Content: "Line 1\nLine 2\tTabbed\rCarriage Return",
				Link:    "http://example.com",
			},
			want: true,
		},
		{
			name: "Invalid - Only newlines in title",
			post: Post{
				Title: "\n\n\n",
				Link:  "http://example.com",
			},
			want: false,
		},
		{
			name: "Invalid - Only newlines in link",
			post: Post{
				Title: "Title",
				Link:  "\n\n\n",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.post.IsValid()
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestPostSanitize tests the Sanitize method
func TestPostSanitize(t *testing.T) {
	tests := []struct {
		name      string
		post      Post
		wantTitle string
		wantContent string
		wantLink  string
	}{
		{
			name: "Trim whitespace from all fields",
			post: Post{
				Title:   "  Title  ",
				Content: "  Content  ",
				Link:    "  http://example.com  ",
			},
			wantTitle:   "Title",
			wantContent: "Content",
			wantLink:    "http://example.com",
		},
		{
			name: "No whitespace to trim",
			post: Post{
				Title:   "Title",
				Content: "Content",
				Link:    "http://example.com",
			},
			wantTitle:   "Title",
			wantContent: "Content",
			wantLink:    "http://example.com",
		},
		{
			name: "Trim tabs and newlines",
			post: Post{
				Title:   "\t\nTitle\n\t",
				Content: "\nContent\t",
				Link:    "\thttp://example.com\n",
			},
			wantTitle:   "Title",
			wantContent: "Content",
			wantLink:    "http://example.com",
		},
		{
			name: "Empty strings remain empty",
			post: Post{
				Title:   "",
				Content: "",
				Link:    "",
			},
			wantTitle:   "",
			wantContent: "",
			wantLink:    "",
		},
		{
			name: "Whitespace only becomes empty",
			post: Post{
				Title:   "   \t\n   ",
				Content: "   ",
				Link:    "\t\t",
			},
			wantTitle:   "",
			wantContent: "",
			wantLink:    "",
		},
		{
			name: "Preserve ID and Read fields",
			post: Post{
				ID:      42,
				Title:   "  Title  ",
				Content: "  Content  ",
				Link:    "  Link  ",
				Read:    true,
			},
			wantTitle:   "Title",
			wantContent: "Content",
			wantLink:    "Link",
		},
		{
			name: "Very long content trimmed correctly",
			post: Post{
				Title:   "  " + strings.Repeat("A", 10000) + "  ",
				Content: "  " + strings.Repeat("Lorem ipsum", 1000) + "  ",
				Link:    "  http://example.com  ",
			},
			wantTitle:   strings.Repeat("A", 10000),
			wantContent: strings.Repeat("Lorem ipsum", 1000),
			wantLink:    "http://example.com",
		},
		{
			name: "Unicode preserved correctly",
			post: Post{
				Title:   "  日本語のタイトル 🚀  ",
				Content: "  محتوى عربي  ",
				Link:    "  http://example.com/日本語  ",
			},
			wantTitle:   "日本語のタイトル 🚀",
			wantContent: "محتوى عربي",
			wantLink:    "http://example.com/日本語",
		},
		{
			name: "HTML content preserved",
			post: Post{
				Title:   "  HTML Post  ",
				Content: "  <script>alert('xss')</script><p>Paragraph</p>  ",
				Link:    "  http://example.com  ",
			},
			wantTitle:   "HTML Post",
			wantContent: "<script>alert('xss')</script><p>Paragraph</p>",
			wantLink:    "http://example.com",
		},
		{
			name: "Newlines and tabs in content preserved (only outer whitespace trimmed)",
			post: Post{
				Title:   "  Multi-line Title\nLine 2  ",
				Content: "  Line 1\nLine 2\tTabbed\rCarriage Return  ",
				Link:    "  http://example.com  ",
			},
			wantTitle:   "Multi-line Title\nLine 2",
			wantContent: "Line 1\nLine 2\tTabbed\rCarriage Return",
			wantLink:    "http://example.com",
		},
		{
			name: "Negative ID preserved",
			post: Post{
				ID:      -1,
				Title:   "  Negative  ",
				Content: "  Content  ",
				Link:    "  http://example.com  ",
			},
			wantTitle:   "Negative",
			wantContent: "Content",
			wantLink:    "http://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.post.Sanitize()
			
			if got.Title != tt.wantTitle {
				t.Errorf("Sanitize().Title = %q, want %q", got.Title, tt.wantTitle)
			}
			if got.Content != tt.wantContent {
				t.Errorf("Sanitize().Content = %q, want %q", got.Content, tt.wantContent)
			}
			if got.Link != tt.wantLink {
				t.Errorf("Sanitize().Link = %q, want %q", got.Link, tt.wantLink)
			}
			
			// Verify ID and Read are preserved
			if got.ID != tt.post.ID {
				t.Errorf("Sanitize().ID = %d, want %d", got.ID, tt.post.ID)
			}
			if got.Read != tt.post.Read {
				t.Errorf("Sanitize().Read = %v, want %v", got.Read, tt.post.Read)
			}
		})
	}
}

// TestSanitizeDoesNotModifyOriginal verifies Sanitize returns a new Post
func TestSanitizeDoesNotModifyOriginal(t *testing.T) {
	original := Post{
		Title:   "  Original  ",
		Content: "  Content  ",
		Link:    "  Link  ",
	}
	
	sanitized := original.Sanitize()
	
	// Original should remain unchanged
	if original.Title != "  Original  " {
		t.Errorf("Sanitize modified original.Title")
	}
	if original.Content != "  Content  " {
		t.Errorf("Sanitize modified original.Content")
	}
	if original.Link != "  Link  " {
		t.Errorf("Sanitize modified original.Link")
	}
	
	// Sanitized should be trimmed
	if sanitized.Title != "Original" {
		t.Errorf("Sanitized.Title = %q, want %q", sanitized.Title, "Original")
	}
}

// TestPostStructFields tests Post struct field values
func TestPostStructFields(t *testing.T) {
	tests := []struct {
		name    string
		post    Post
		wantID  int
		wantRead bool
	}{
		{
			name: "Default Post values",
			post: Post{
				Title:   "Test Post",
				Content: "Test Content",
				Link:    "http://example.com",
			},
			wantID:   0,
			wantRead: false,
		},
		{
			name: "Post with all fields set",
			post: Post{
				ID:      42,
				Title:   "Full Post",
				Content: "Full Content",
				Link:    "http://example.com/full",
				Read:    true,
			},
			wantID:   42,
			wantRead: true,
		},
		{
			name: "Post with empty strings",
			post: Post{
				Title:   "",
				Content: "",
				Link:    "",
			},
			wantID:   0,
			wantRead: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.post.ID != tt.wantID {
				t.Errorf("ID = %d, want %d", tt.post.ID, tt.wantID)
			}
			if tt.post.Read != tt.wantRead {
				t.Errorf("Read = %v, want %v", tt.post.Read, tt.wantRead)
			}
		})
	}
}

// TestFeedStructFields tests Feed struct field values
func TestFeedStructFields(t *testing.T) {
	tests := []struct {
		name      string
		feed      Feed
		wantID    int
		wantPosts int
	}{
		{
			name: "Empty Feed",
			feed: Feed{
				Title: "Test Feed",
				URL:   "http://example.com",
			},
			wantID:    0,
			wantPosts: 0,
		},
		{
			name: "Feed with posts",
			feed: Feed{
				ID:    1,
				Title: "Test Feed",
				URL:   "http://example.com",
				Posts: []Post{
					{Title: "Post 1"},
					{Title: "Post 2"},
				},
			},
			wantID:    1,
			wantPosts: 2,
		},
		{
			name: "Feed with nil posts slice",
			feed: Feed{
				Title: "Test Feed",
				URL:   "http://example.com",
				Posts: nil,
			},
			wantID:    0,
			wantPosts: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.feed.ID != tt.wantID {
				t.Errorf("ID = %d, want %d", tt.feed.ID, tt.wantID)
			}
			if len(tt.feed.Posts) != tt.wantPosts {
				t.Errorf("Posts count = %d, want %d", len(tt.feed.Posts), tt.wantPosts)
			}
		})
	}
}

// TestFeedEdgeCases tests Feed with edge case values  
func TestFeedEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		feed        Feed
		description string
	}{
		{
			name: "Feed with many posts",
			feed: Feed{
				Title: "Large Feed",
				URL:   "http://example.com",
				Posts: make([]Post, 1000),
			},
			description: "Feed should handle large number of posts",
		},
		{
			name: "Feed with empty URL",
			feed: Feed{
				Title: "No URL",
				URL:   "",
				Posts: []Post{{Title: "Post"}},
			},
			description: "Feed can have empty URL",
		},
		{
			name: "Feed with very long URL",
			feed: Feed{
				Title: "Long URL",
				URL:   "http://example.com/" + strings.Repeat("path/", 500),
				Posts: nil,
			},
			description: "Feed should handle very long URLs",
		},
		{
			name: "Feed with unicode in URL",
			feed: Feed{
				Title: "Unicode URL",
				URL:   "http://example.com/日本語/feed.xml",
			},
			description: "Feed should handle unicode in URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify struct can be created and fields accessible
			_ = tt.feed.ID
			_ = tt.feed.Title
			_ = tt.feed.URL
			_ = tt.feed.Posts
			
			t.Logf("Test case: %s", tt.description)
		})
	}
}

// TestPostComparison tests comparing Post structs
func TestPostComparison(t *testing.T) {
	post1 := Post{
		ID:      1,
		Title:   "Same Title",
		Content: "Same Content",
		Link:    "http://example.com",
		Read:    false,
	}
	
	post2 := Post{
		ID:      1,
		Title:   "Same Title",
		Content: "Same Content",
		Link:    "http://example.com",
		Read:    false,
	}
	
	post3 := Post{
		ID:      2,
		Title:   "Different Title",
		Content: "Different Content",
		Link:    "http://example.com/different",
		Read:    true,
	}
	
	// Test field-by-field comparison
	if post1.ID != post2.ID {
		t.Errorf("post1.ID != post2.ID")
	}
	if post1.Title != post2.Title {
		t.Errorf("post1.Title != post2.Title")
	}
	
	if post1.ID == post3.ID {
		t.Errorf("post1.ID should not equal post3.ID")
	}
	if post1.Read == post3.Read {
		t.Errorf("post1.Read should not equal post3.Read")
	}
}
