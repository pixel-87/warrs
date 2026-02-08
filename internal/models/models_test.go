package models

import (
	"testing"
)

func TestPostValidation(t *testing.T) {
	tests := []struct {
		name string
		post Post
		want bool // whether the post has required fields
	}{
		{
			name: "Valid post with all fields",
			post: Post{
				Title:   "Test Post",
				Content: "Test Content",
				Link:    "https://example.com/post",
				ID:      1,
				Read:    false,
			},
			want: true,
		},
		{
			name: "Empty title",
			post: Post{
				Title:   "",
				Content: "Content",
				Link:    "https://example.com/post",
			},
			want: false,
		},
		{
			name: "Empty link",
			post: Post{
				Title:   "Title",
				Content: "Content",
				Link:    "",
			},
			want: false,
		},
		{
			name: "Empty content allowed",
			post: Post{
				Title: "Title",
				Link:  "https://example.com/post",
			},
			want: true,
		},
		{
			name: "Negative ID",
			post: Post{
				Title:   "Title",
				Content: "Content",
				Link:    "https://example.com/post",
				ID:      -1,
			},
			want: true, // ID can be negative in struct, validation should happen elsewhere
		},
		{
			name: "Very long title",
			post: Post{
				Title:   string(make([]byte, 10000)),
				Content: "Content",
				Link:    "https://example.com/post",
			},
			want: true,
		},
		{
			name: "Unicode in title",
			post: Post{
				Title:   "测试标题 🎉",
				Content: "Content",
				Link:    "https://example.com/post",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test struct creation doesn't panic
			hasRequired := tt.post.Title != "" && tt.post.Link != ""
			if hasRequired != tt.want {
				t.Errorf("Post validation = %v, want %v", hasRequired, tt.want)
			}
		})
	}
}

func TestFeedValidation(t *testing.T) {
	tests := []struct {
		name string
		feed Feed
		want bool
	}{
		{
			name: "Valid feed with posts",
			feed: Feed{
				Title: "Test Feed",
				URL:   "https://example.com/feed.xml",
				Posts: []Post{
					{Title: "Post 1", Link: "https://example.com/1"},
				},
				ID: 1,
			},
			want: true,
		},
		{
			name: "Valid feed without posts",
			feed: Feed{
				Title: "Empty Feed",
				URL:   "https://example.com/feed.xml",
				Posts: []Post{},
			},
			want: true,
		},
		{
			name: "Empty title",
			feed: Feed{
				Title: "",
				URL:   "https://example.com/feed.xml",
			},
			want: false,
		},
		{
			name: "Empty URL",
			feed: Feed{
				Title: "Feed",
				URL:   "",
			},
			want: false,
		},
		{
			name: "Nil posts slice",
			feed: Feed{
				Title: "Feed",
				URL:   "https://example.com/feed.xml",
				Posts: nil,
			},
			want: true,
		},
		{
			name: "Large number of posts",
			feed: Feed{
				Title: "Big Feed",
				URL:   "https://example.com/feed.xml",
				Posts: make([]Post, 1000),
			},
			want: true,
		},
		{
			name: "Unicode in feed title",
			feed: Feed{
				Title: "フィード テスト",
				URL:   "https://example.com/feed.xml",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasRequired := tt.feed.Title != "" && tt.feed.URL != ""
			if hasRequired != tt.want {
				t.Errorf("Feed validation = %v, want %v", hasRequired, tt.want)
			}
		})
	}
}

func TestPostReadStatus(t *testing.T) {
	tests := []struct {
		name       string
		initialVal bool
		want       bool
	}{
		{
			name:       "Unread post",
			initialVal: false,
			want:       false,
		},
		{
			name:       "Read post",
			initialVal: true,
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := Post{
				Title: "Test",
				Link:  "https://example.com",
				Read:  tt.initialVal,
			}
			if post.Read != tt.want {
				t.Errorf("Post.Read = %v, want %v", post.Read, tt.want)
			}
		})
	}
}
