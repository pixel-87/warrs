package rss

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr error
	}{
		// Happy path
		{
			name:    "valid http URL",
			url:     "http://example.com/feed.xml",
			wantErr: nil,
		},
		{
			name:    "valid https URL",
			url:     "https://example.com/feed.xml",
			wantErr: nil,
		},
		{
			name:    "URL with path and query",
			url:     "https://blog.example.com/rss?format=xml",
			wantErr: nil,
		},
		{
			name:    "URL with port",
			url:     "http://localhost:8080/feed",
			wantErr: nil,
		},
		{
			name:    "URL with subdomain",
			url:     "https://feeds.example.co.uk/rss.xml",
			wantErr: nil,
		},
		
		// Edge cases - Empty/Whitespace
		{
			name:    "empty URL",
			url:     "",
			wantErr: ErrEmptyURL,
		},
		{
			name:    "whitespace only URL",
			url:     "   ",
			wantErr: ErrEmptyURL,
		},
		{
			name:    "tab only URL",
			url:     "\t\t",
			wantErr: ErrEmptyURL,
		},
		{
			name:    "newline only URL",
			url:     "\n",
			wantErr: ErrEmptyURL,
		},
		
		// Edge cases - Invalid schemes
		{
			name:    "ftp scheme",
			url:     "ftp://example.com/feed",
			wantErr: ErrInvalidScheme,
		},
		{
			name:    "file scheme",
			url:     "file:///tmp/feed.xml",
			wantErr: ErrInvalidScheme,
		},
		{
			name:    "no scheme",
			url:     "example.com/feed",
			wantErr: ErrInvalidScheme,
		},
		{
			name:    "javascript scheme",
			url:     "javascript:alert(1)",
			wantErr: ErrInvalidScheme,
		},
		{
			name:    "data URI",
			url:     "data:text/plain,hello",
			wantErr: ErrInvalidScheme,
		},
		
		// Edge cases - Malformed URLs
		{
			name:    "missing host",
			url:     "http://",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "invalid characters",
			url:     "http://exam ple.com/feed",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "only scheme and slashes",
			url:     "https://",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "malformed with colons",
			url:     "http:::example.com",
			wantErr: ErrInvalidURL,
		},
		
		// Edge cases - Special characters
		{
			name:    "URL with fragment",
			url:     "https://example.com/feed#top",
			wantErr: nil,
		},
		{
			name:    "URL with encoded characters",
			url:     "https://example.com/feed%20rss.xml",
			wantErr: nil,
		},
		{
			name:    "URL with international domain",
			url:     "https://münchen.de/feed",
			wantErr: nil,
		},
		{
			name:    "URL with IP address",
			url:     "http://192.168.1.1/feed",
			wantErr: nil,
		},
		{
			name:    "URL with IPv6",
			url:     "http://[::1]/feed",
			wantErr: nil,
		},
		
		// Edge cases - Very long URLs
		{
			name:    "very long URL",
			url:     "https://example.com/" + strings.Repeat("a", 2000) + "/feed.xml",
			wantErr: nil,
		},
		
		// Edge cases - URL with credentials (should be valid but not recommended)
		{
			name:    "URL with username and password",
			url:     "https://user:pass@example.com/feed",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ValidateURL() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateURL() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr error
	}{
		// Happy path
		{
			name:    "valid simple title",
			title:   "My Blog",
			wantErr: nil,
		},
		{
			name:    "valid title with spaces",
			title:   "Tech News and Updates",
			wantErr: nil,
		},
		{
			name:    "valid title with special characters",
			title:   "Bob's Tech Blog: 2024 Edition!",
			wantErr: nil,
		},
		{
			name:    "valid title with unicode",
			title:   "日本語ブログ",
			wantErr: nil,
		},
		{
			name:    "valid title with emojis",
			title:   "🚀 Space News 🌟",
			wantErr: nil,
		},
		
		// Edge cases - Empty/Whitespace
		{
			name:    "empty title",
			title:   "",
			wantErr: ErrEmptyTitle,
		},
		{
			name:    "whitespace only title",
			title:   "   ",
			wantErr: ErrEmptyTitle,
		},
		{
			name:    "tab only title",
			title:   "\t\t\t",
			wantErr: ErrEmptyTitle,
		},
		{
			name:    "newline only title",
			title:   "\n\n",
			wantErr: ErrEmptyTitle,
		},
		{
			name:    "mixed whitespace title",
			title:   " \t\n\r ",
			wantErr: ErrEmptyTitle,
		},
		
		// Edge cases - Length
		{
			name:    "title at max length",
			title:   strings.Repeat("a", MaxTitleLength),
			wantErr: nil,
		},
		{
			name:    "title one char over max",
			title:   strings.Repeat("a", MaxTitleLength+1),
			wantErr: ErrTitleTooLong,
		},
		{
			name:    "title way over max length",
			title:   strings.Repeat("a", MaxTitleLength*10),
			wantErr: ErrTitleTooLong,
		},
		
		// Edge cases - Whitespace handling
		{
			name:    "title with leading whitespace",
			title:   "   Leading Spaces",
			wantErr: nil,
		},
		{
			name:    "title with trailing whitespace",
			title:   "Trailing Spaces   ",
			wantErr: nil,
		},
		{
			name:    "title with both leading and trailing whitespace",
			title:   "  Both Sides  ",
			wantErr: nil,
		},
		
		// Edge cases - Single character
		{
			name:    "single character title",
			title:   "A",
			wantErr: nil,
		},
		{
			name:    "single unicode character",
			title:   "あ",
			wantErr: nil,
		},
		{
			name:    "single emoji",
			title:   "🎉",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTitle(tt.title)
			
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ValidateTitle() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("ValidateTitle() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateTitle() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestSanitizeTitle(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  string
	}{
		// Happy path
		{
			name:  "normal title unchanged",
			title: "My Blog",
			want:  "My Blog",
		},
		{
			name:  "title with spaces unchanged",
			title: "Tech News Today",
			want:  "Tech News Today",
		},
		
		// Edge cases - Whitespace normalization
		{
			name:  "leading spaces removed",
			title: "   Leading",
			want:  "Leading",
		},
		{
			name:  "trailing spaces removed",
			title: "Trailing   ",
			want:  "Trailing",
		},
		{
			name:  "both leading and trailing spaces removed",
			title: "  Both  ",
			want:  "Both",
		},
		{
			name:  "multiple spaces collapsed",
			title: "Multiple    Spaces    Here",
			want:  "Multiple Spaces Here",
		},
		{
			name:  "tabs replaced with spaces",
			title: "Tab\tSeparated\tWords",
			want:  "Tab Separated Words",
		},
		{
			name:  "newlines replaced with spaces",
			title: "Line\nBreak\nHere",
			want:  "Line Break Here",
		},
		{
			name:  "carriage returns replaced",
			title: "Carriage\rReturn\rHere",
			want:  "Carriage Return Here",
		},
		{
			name:  "mixed whitespace normalized",
			title: "Mixed \t\n\r  Whitespace",
			want:  "Mixed Whitespace",
		},
		
		// Edge cases - Empty/Whitespace input
		{
			name:  "empty string",
			title: "",
			want:  "",
		},
		{
			name:  "only spaces",
			title: "     ",
			want:  "",
		},
		{
			name:  "only tabs",
			title: "\t\t\t",
			want:  "",
		},
		{
			name:  "only newlines",
			title: "\n\n\n",
			want:  "",
		},
		
		// Edge cases - Length truncation
		{
			name:  "title at max length unchanged",
			title: strings.Repeat("a", MaxTitleLength),
			want:  strings.Repeat("a", MaxTitleLength),
		},
		{
			name:  "title over max length truncated",
			title: strings.Repeat("b", MaxTitleLength+100),
			want:  strings.Repeat("b", MaxTitleLength),
		},
		{
			name:  "title with spaces over max truncated",
			title: strings.Repeat("c ", MaxTitleLength),
			want:  strings.Repeat("c ", MaxTitleLength/2)[:MaxTitleLength],
		},
		
		// Edge cases - Special characters preserved
		{
			name:  "unicode characters preserved",
			title: "日本語 タイトル",
			want:  "日本語 タイトル",
		},
		{
			name:  "emojis preserved",
			title: "🎉 Party Time 🎊",
			want:  "🎉 Party Time 🎊",
		},
		{
			name:  "special punctuation preserved",
			title: "Bob's Blog: The \"Best\" One!",
			want:  "Bob's Blog: The \"Best\" One!",
		},
		
		// Edge cases - Complex scenarios
		{
			name:  "multiple issues combined",
			title: "  \tMultiple\n\n  Spaces  \r\n  And    Breaks  ",
			want:  "Multiple Spaces And Breaks",
		},
		{
			name:  "single character",
			title: "A",
			want:  "A",
		},
		{
			name:  "single character with whitespace",
			title: "  A  ",
			want:  "A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeTitle(tt.title)
			if got != tt.want {
				t.Errorf("SanitizeTitle() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		// Happy path
		{
			name:    "simple URL unchanged",
			url:     "https://example.com/feed",
			want:    "https://example.com/feed",
			wantErr: false,
		},
		{
			name:    "URL with path",
			url:     "https://blog.example.com/rss.xml",
			want:    "https://blog.example.com/rss.xml",
			wantErr: false,
		},
		
		// Edge cases - Fragment removal
		{
			name:    "fragment removed",
			url:     "https://example.com/feed#section",
			want:    "https://example.com/feed",
			wantErr: false,
		},
		{
			name:    "multiple fragments removed",
			url:     "https://example.com/feed#top#bottom",
			want:    "https://example.com/feed",
			wantErr: false,
		},
		
		// Edge cases - Case normalization
		{
			name:    "uppercase scheme lowercased",
			url:     "HTTPS://example.com/feed",
			want:    "https://example.com/feed",
			wantErr: false,
		},
		{
			name:    "uppercase host lowercased",
			url:     "https://EXAMPLE.COM/feed",
			want:    "https://example.com/feed",
			wantErr: false,
		},
		{
			name:    "mixed case scheme and host lowercased",
			url:     "HtTpS://ExAmPlE.CoM/feed",
			want:    "https://example.com/feed",
			wantErr: false,
		},
		{
			name:    "path case preserved",
			url:     "https://example.com/Feed/RSS.XML",
			want:    "https://example.com/Feed/RSS.XML",
			wantErr: false,
		},
		
		// Edge cases - Whitespace
		{
			name:    "leading whitespace trimmed",
			url:     "  https://example.com/feed",
			want:    "https://example.com/feed",
			wantErr: false,
		},
		{
			name:    "trailing whitespace trimmed",
			url:     "https://example.com/feed  ",
			want:    "https://example.com/feed",
			wantErr: false,
		},
		{
			name:    "both whitespace trimmed",
			url:     "  https://example.com/feed  ",
			want:    "https://example.com/feed",
			wantErr: false,
		},
		
		// Edge cases - Query parameters preserved
		{
			name:    "query params preserved",
			url:     "https://example.com/feed?format=xml&limit=10",
			want:    "https://example.com/feed?format=xml&limit=10",
			wantErr: false,
		},
		{
			name:    "query params with fragment removed",
			url:     "https://example.com/feed?format=xml#top",
			want:    "https://example.com/feed?format=xml",
			wantErr: false,
		},
		
		// Edge cases - Invalid URLs return errors
		{
			name:    "empty URL returns error",
			url:     "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid scheme returns error",
			url:     "ftp://example.com/feed",
			want:    "",
			wantErr: true,
		},
		{
			name:    "no host returns error",
			url:     "https://",
			want:    "",
			wantErr: true,
		},
		
		// Edge cases - Port preserved
		{
			name:    "port number preserved",
			url:     "http://localhost:8080/feed",
			want:    "http://localhost:8080/feed",
			wantErr: false,
		},
		{
			name:    "non-standard port preserved",
			url:     "https://example.com:9443/feed",
			want:    "https://example.com:9443/feed",
			wantErr: false,
		},
		
		// Edge cases - Special characters
		{
			name:    "encoded characters preserved",
			url:     "https://example.com/feed%20rss.xml",
			want:    "https://example.com/feed%20rss.xml",
			wantErr: false,
		},
		{
			name:    "IP address normalized",
			url:     "HTTP://192.168.1.1/feed",
			want:    "http://192.168.1.1/feed",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeURL(tt.url)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("NormalizeURL() error = nil, wantErr = true")
				}
				return
			}
			
			if err != nil {
				t.Errorf("NormalizeURL() unexpected error = %v", err)
				return
			}
			
			if got != tt.want {
				t.Errorf("NormalizeURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsValidXMLStart(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		// Happy path - Valid XML starts
		{
			name:    "standard XML declaration",
			content: `<?xml version="1.0" encoding="UTF-8"?>`,
			want:    true,
		},
		{
			name:    "RSS feed start",
			content: `<rss version="2.0">`,
			want:    true,
		},
		{
			name:    "Atom feed start",
			content: `<feed xmlns="http://www.w3.org/2005/Atom">`,
			want:    true,
		},
		{
			name:    "XML with leading whitespace",
			content: `   <?xml version="1.0"?>`,
			want:    true,
		},
		{
			name:    "RSS with newlines",
			content: "\n\n<rss version=\"2.0\">",
			want:    true,
		},
		
		// Edge cases - Empty/Whitespace
		{
			name:    "empty string",
			content: "",
			want:    false,
		},
		{
			name:    "only spaces",
			content: "     ",
			want:    false,
		},
		{
			name:    "only tabs",
			content: "\t\t\t",
			want:    false,
		},
		{
			name:    "only newlines",
			content: "\n\n\n",
			want:    false,
		},
		{
			name:    "mixed whitespace",
			content: " \t\n\r ",
			want:    false,
		},
		
		// Edge cases - Invalid starts
		{
			name:    "HTML document",
			content: `<!DOCTYPE html><html>`,
			want:    false,
		},
		{
			name:    "plain text",
			content: "This is just plain text",
			want:    false,
		},
		{
			name:    "JSON content",
			content: `{"feed": "data"}`,
			want:    false,
		},
		{
			name:    "partial XML declaration",
			content: `<?xml`,
			want:    true, // starts with <?xml
		},
		{
			name:    "partial RSS tag",
			content: `<rss`,
			want:    true, // starts with <rss
		},
		{
			name:    "partial feed tag",
			content: `<feed`,
			want:    true, // starts with <feed
		},
		{
			name:    "RSS with wrong case",
			content: `<RSS version="2.0">`,
			want:    false, // case sensitive
		},
		{
			name:    "Feed with wrong case",
			content: `<FEED>`,
			want:    false, // case sensitive
		},
		
		// Edge cases - Almost valid
		{
			name:    "XML with typo",
			content: `<?xm version="1.0"?>`,
			want:    false,
		},
		{
			name:    "RSS with extra character",
			content: `x<rss version="2.0">`,
			want:    false,
		},
		{
			name:    "comment before XML",
			content: `<!-- comment --><?xml version="1.0"?>`,
			want:    false,
		},
		
		// Edge cases - Unusual but valid
		{
			name:    "XML declaration with lots of whitespace",
			content: "  \t\n  <?xml version=\"1.0\"?>",
			want:    true,
		},
		{
			name:    "RSS 1.0 format",
			content: `<rss version="1.0">`,
			want:    true,
		},
		
		// Edge cases - Binary/garbage data
		{
			name:    "binary data",
			content: "\x00\x01\x02\x03",
			want:    false,
		},
		{
			name:    "unicode BOM with XML",
			content: "\uFEFF<?xml version=\"1.0\"?>",
			want:    false, // BOM is not trimmed by TrimSpace
		},
		
		// Edge cases - Very long content
		{
			name:    "valid XML with long content",
			content: "<?xml version=\"1.0\"?>" + strings.Repeat("<item>test</item>", 1000),
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidXMLStart(tt.content)
			if got != tt.want {
				t.Errorf("IsValidXMLStart() = %v, want %v", got, tt.want)
			}
		})
	}
}
