package rss

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFeed(t *testing.T) {
	tests := []struct {
		name      string
		filename  string // path to testdata
		wantTitle string
		wantCount int
	}{
		{
			name:      "Basic RSS Feed",
			filename:  "test_feed.xml",
			wantTitle: "Test Blog",
			wantCount: 1,
		},

		{
			name:      "Content Fallback test",
			filename:  "test_fallback.xml",
			wantTitle: "Fallback Test",
			wantCount: 1,
		},
	}

	f := NewFetcher()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join("testdata", tt.filename)

			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("couldn't read test file %s: %v", tt.filename, err)
			}

			res, err := f.parseFeed("http://test.com", content)
			if err != nil {
				t.Fatalf("parseFeed failed: %v", err)
			}

			if res.Title != tt.wantTitle {
				t.Errorf("got title %q, want %q", res.Title, tt.wantTitle)
			}

			if len(res.Posts) != tt.wantCount {
				t.Errorf("got %d posts, want %d", len(res.Posts), tt.wantCount)
			}
		})
	}
}
