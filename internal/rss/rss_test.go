package rss

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFeed(t *testing.T) {
	tests := []struct {
		name      string
		filename  string // path to testdata
		wantTitle string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "Basic RSS Feed",
			filename:  "test_feed.xml",
			wantTitle: "Test Blog",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "Content Fallback test",
			filename:  "test_fallback.xml",
			wantTitle: "Fallback Test",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "Empty Feed - No Items",
			filename:  "empty_feed.xml",
			wantTitle: "Empty Feed",
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "Feed Without Title",
			filename:  "no_title_feed.xml",
			wantTitle: "",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "Unicode and Emoji Content",
			filename:  "unicode_feed.xml",
			wantTitle: "Unicode Test Blog 🚀 日本語 العربية",
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "Malformed XML",
			filename:  "malformed.xml",
			wantTitle: "",
			wantCount: 0,
			wantErr:   true,
		},
		{
			name:      "Empty Strings Feed",
			filename:  "empty_strings.xml",
			wantTitle: "",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "Large Feed with Multiple Items",
			filename:  "large_feed.xml",
			wantTitle: "Large Feed",
			wantCount: 5,
			wantErr:   false,
		},
	}

	f := NewFetcher(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join("testdata", tt.filename)

			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("couldn't read test file %s: %v", tt.filename, err)
			}

			res, err := f.parseFeed("http://test.com", content)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseFeed() expected error but got none")
				}
				return
			}
			
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

// TestParseFeedEdgeCases tests edge cases with inline data
func TestParseFeedEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		url       string
		wantTitle string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "Completely Empty Data",
			data:      "",
			url:       "http://test.com",
			wantTitle: "",
			wantCount: 0,
			wantErr:   true,
		},
		{
			name:      "Invalid XML - Not RSS",
			data:      "<html><body>Not an RSS feed</body></html>",
			url:       "http://test.com",
			wantTitle: "",
			wantCount: 0,
			wantErr:   true,
		},
		{
			name:      "Null Bytes in Data",
			data:      "<?xml version=\"1.0\"?>\x00<rss><channel><title>Bad</title></channel></rss>",
			url:       "http://test.com",
			wantTitle: "",
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "Very Long Title",
			data: `<?xml version="1.0"?><rss version="2.0"><channel><title>` + strings.Repeat("A", 1000) + `</title></channel></rss>`,
			url:  "http://test.com",
			wantTitle: strings.Repeat("A", 1000),
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "HTML Entities in Title",
			data: `<?xml version="1.0"?><rss version="2.0"><channel><title>&lt;Test&gt; &amp; Blog</title></channel></rss>`,
			url:  "http://test.com",
			wantTitle: "<Test> & Blog",
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "Missing Both Content and Description",
			data: `<?xml version="1.0"?><rss version="2.0"><channel><title>Test</title><item><title>Post</title><link>http://example.com</link></item></channel></rss>`,
			url:  "http://test.com",
			wantTitle: "Test",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "Just Whitespace",
			data:      "   \n\t  \r\n  ",
			url:       "http://test.com",
			wantTitle: "",
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "CDATA Content",
			data: `<?xml version="1.0"?><rss version="2.0"><channel><title>Test</title><item><title>Post</title><link>http://example.com</link><description><![CDATA[<p>HTML content</p>]]></description></item></channel></rss>`,
			url:  "http://test.com",
			wantTitle: "Test",
			wantCount: 1,
			wantErr:   false,
		},
	}

	f := NewFetcher(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := f.parseFeed(tt.url, []byte(tt.data))
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseFeed() expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Fatalf("parseFeed() unexpected error: %v", err)
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

// TestParseFeedContentPriority verifies content vs description priority
func TestParseFeedContentPriority(t *testing.T) {
	tests := []struct {
		name        string
		data        string
		wantContent string
	}{
		{
			name: "Content Takes Priority Over Description",
			data: `<?xml version="1.0"?><rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/"><channel><title>Test</title><item><title>Post</title><link>http://example.com</link><content:encoded>Priority Content</content:encoded><description>Fallback Description</description></item></channel></rss>`,
			wantContent: "Priority Content",
		},
		{
			name: "Description Used When No Content",
			data: `<?xml version="1.0"?><rss version="2.0"><channel><title>Test</title><item><title>Post</title><link>http://example.com</link><description>Fallback Description</description></item></channel></rss>`,
			wantContent: "Fallback Description",
		},
		{
			name: "Empty String When Both Missing",
			data: `<?xml version="1.0"?><rss version="2.0"><channel><title>Test</title><item><title>Post</title><link>http://example.com</link></item></channel></rss>`,
			wantContent: "",
		},
		{
			name: "Empty Content Triggers Fallback to Description",
			data: `<?xml version="1.0"?><rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/"><channel><title>Test</title><item><title>Post</title><link>http://example.com</link><content:encoded></content:encoded><description>Should Use This</description></item></channel></rss>`,
			wantContent: "Should Use This",
		},
	}

	f := NewFetcher(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := f.parseFeed("http://test.com", []byte(tt.data))
			if err != nil {
				t.Fatalf("parseFeed() unexpected error: %v", err)
			}

			if len(res.Posts) == 0 {
				t.Fatal("expected at least one post")
			}

			if res.Posts[0].Content != tt.wantContent {
				t.Errorf("got content %q, want %q", res.Posts[0].Content, tt.wantContent)
			}
		})
	}
}
