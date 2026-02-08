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
			name:      "Empty feed - no items",
			filename:  "empty_feed.xml",
			wantTitle: "Empty Feed",
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "Unicode characters in feed",
			filename:  "unicode_feed.xml",
			wantTitle: "Unicode 测试 フィード 🎉",
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "Large feed with multiple items",
			filename:  "large_feed.xml",
			wantTitle: "Large Feed",
			wantCount: 5,
			wantErr:   false,
		},
		{
			name:      "Feed with empty items",
			filename:  "empty_items.xml",
			wantTitle: "Feed with Empty Items",
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
					t.Errorf("parseFeed() expected error, got nil")
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

func TestParseFeedEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		data    []byte
		wantErr bool
	}{
		{
			name:    "Empty data",
			url:     "http://test.com",
			data:    []byte{},
			wantErr: true,
		},
		{
			name:    "Nil data",
			url:     "http://test.com",
			data:    nil,
			wantErr: true,
		},
		{
			name:    "Non-XML data",
			url:     "http://test.com",
			data:    []byte("this is not xml"),
			wantErr: true,
		},
		{
			name:    "JSON instead of XML",
			url:     "http://test.com",
			data:    []byte(`{"feed": "json"}`),
			wantErr: true,
		},
		{
			name:    "Empty URL with valid data",
			url:     "",
			data:    []byte(`<?xml version="1.0"?><rss version="2.0"><channel><title>Test</title></channel></rss>`),
			wantErr: false, // URL is just stored, not validated
		},
		{
			name:    "Very long URL",
			url:     "http://example.com/" + strings.Repeat("a", 10000),
			data:    []byte(`<?xml version="1.0"?><rss version="2.0"><channel><title>Test</title></channel></rss>`),
			wantErr: false,
		},
		{
			name:    "Special characters in URL",
			url:     "http://example.com/feed?param=测试&emoji=🎉",
			data:    []byte(`<?xml version="1.0"?><rss version="2.0"><channel><title>Test</title></channel></rss>`),
			wantErr: false,
		},
	}

	f := NewFetcher(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := f.parseFeed(tt.url, tt.data)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseFeed() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("parseFeed() unexpected error: %v", err)
			}

			if res.URL != tt.url {
				t.Errorf("got URL %q, want %q", res.URL, tt.url)
			}
		})
	}
}

func TestParseFeedContentPriority(t *testing.T) {
	// Test that Content field is preferred over Description
	xmlWithContent := `<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/">
<channel>
    <title>Content Priority Test</title>
    <item>
        <title>Post with Content</title>
        <link>http://example.com/1</link>
        <description>Description text</description>
        <content:encoded><![CDATA[Full content text]]></content:encoded>
    </item>
</channel>
</rss>`

	f := NewFetcher(nil)
	res, err := f.parseFeed("http://test.com", []byte(xmlWithContent))

	if err != nil {
		t.Fatalf("parseFeed failed: %v", err)
	}

	if len(res.Posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(res.Posts))
	}

	// gofeed should prefer content:encoded over description
	if res.Posts[0].Content == "" {
		t.Errorf("Post content is empty")
	}
}

func TestParseFeedPostFields(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">
<channel>
    <title>Field Test</title>
    <item>
        <title>Test Post Title</title>
        <link>http://example.com/test</link>
        <description>Test Description</description>
    </item>
</channel>
</rss>`

	f := NewFetcher(nil)
	res, err := f.parseFeed("http://test.com/feed", []byte(xml))

	if err != nil {
		t.Fatalf("parseFeed failed: %v", err)
	}

	if len(res.Posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(res.Posts))
	}

	post := res.Posts[0]

	if post.Title != "Test Post Title" {
		t.Errorf("got title %q, want %q", post.Title, "Test Post Title")
	}

	if post.Link != "http://example.com/test" {
		t.Errorf("got link %q, want %q", post.Link, "http://example.com/test")
	}

	if post.Content != "Test Description" {
		t.Errorf("got content %q, want %q", post.Content, "Test Description")
	}

	// ID and Read should have default values
	if post.ID != 0 {
		t.Errorf("expected default ID 0, got %d", post.ID)
	}

	if post.Read != false {
		t.Errorf("expected default Read false, got %v", post.Read)
	}
}
