package rss

import (
	"strings"
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr error
	}{
		// Happy paths
		{
			name:    "valid http URL",
			url:     "http://example.com/feed.rss",
			wantErr: nil,
		},
		{
			name:    "valid https URL",
			url:     "https://example.com/feed.rss",
			wantErr: nil,
		},
		{
			name:    "valid URL with port",
			url:     "https://example.com:8080/feed",
			wantErr: nil,
		},
		{
			name:    "valid URL with query params",
			url:     "https://example.com/feed?format=rss",
			wantErr: nil,
		},
		{
			name:    "valid URL with path and fragment",
			url:     "https://example.com/blog/feed.xml#section",
			wantErr: nil,
		},
		
		// Edge cases - empty/whitespace
		{
			name:    "empty string",
			url:     "",
			wantErr: ErrEmptyURL,
		},
		{
			name:    "only spaces",
			url:     "   ",
			wantErr: ErrEmptyURL,
		},
		{
			name:    "only tabs",
			url:     "\t\t",
			wantErr: ErrEmptyURL,
		},
		{
			name:    "only newlines",
			url:     "\n\n",
			wantErr: ErrEmptyURL,
		},
		{
			name:    "mixed whitespace",
			url:     " \t\n ",
			wantErr: ErrEmptyURL,
		},
		
		// Edge cases - invalid schemes
		{
			name:    "ftp scheme",
			url:     "ftp://example.com/feed",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "file scheme",
			url:     "file:///path/to/feed.xml",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "javascript scheme",
			url:     "javascript:alert(1)",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "data scheme",
			url:     "data:text/html,<script>alert(1)</script>",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "no scheme",
			url:     "example.com/feed",
			wantErr: ErrInvalidURL,
		},
		
		// Edge cases - malformed URLs
		{
			name:    "missing host",
			url:     "http:///feed",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "only scheme",
			url:     "https://",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "invalid characters",
			url:     "http://exam ple.com/feed",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "multiple schemes",
			url:     "http://https://example.com",
			wantErr: nil, // url.Parse accepts this as valid with host "https:"
		},
		
		// Edge cases - special characters
		{
			name:    "URL with unicode",
			url:     "https://例え.jp/feed",
			wantErr: nil,
		},
		{
			name:    "URL with encoded characters",
			url:     "https://example.com/feed%20url",
			wantErr: nil,
		},
		
		// Edge cases - extreme lengths
		{
			name:    "very long domain",
			url:     "https://" + strings.Repeat("a", 200) + ".com/feed",
			wantErr: nil,
		},
		{
			name:    "very long path",
			url:     "https://example.com/" + strings.Repeat("path/", 100) + "feed.xml",
			wantErr: nil,
		},
		
		// Edge cases - URL with username/password
		{
			name:    "URL with auth",
			url:     "https://user:pass@example.com/feed",
			wantErr: nil,
		},
		
		// Edge cases - localhost and IPs
		{
			name:    "localhost",
			url:     "http://localhost/feed",
			wantErr: nil,
		},
		{
			name:    "IPv4 address",
			url:     "http://192.168.1.1/feed",
			wantErr: nil,
		},
		{
			name:    "IPv6 address",
			url:     "http://[::1]/feed",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if err != tt.wantErr {
				t.Errorf("ValidateURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeTitle(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Happy paths
		{
			name:  "normal title",
			input: "My Blog Post",
			want:  "My Blog Post",
		},
		{
			name:  "title with leading spaces",
			input: "   Title",
			want:  "Title",
		},
		{
			name:  "title with trailing spaces",
			input: "Title   ",
			want:  "Title",
		},
		{
			name:  "title with multiple internal spaces",
			input: "Hello    World",
			want:  "Hello World",
		},
		{
			name:  "title with mixed spacing",
			input: "  Hello    World  ",
			want:  "Hello World",
		},
		
		// Edge cases - empty/whitespace
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "only spaces",
			input: "     ",
			want:  "",
		},
		{
			name:  "only tabs",
			input: "\t\t\t",
			want:  "",
		},
		{
			name:  "only newlines",
			input: "\n\n\n",
			want:  "",
		},
		{
			name:  "mixed whitespace only",
			input: " \t\n\r ",
			want:  "",
		},
		
		// Edge cases - special characters
		{
			name:  "title with unicode",
			input: "  Hello 世界  ",
			want:  "Hello 世界",
		},
		{
			name:  "title with emojis",
			input: "🔥  Hot  News  🔥",
			want:  "🔥 Hot News 🔥",
		},
		{
			name:  "title with punctuation",
			input: "Hello,  World!  How  are  you?",
			want:  "Hello, World! How are you?",
		},
		
		// Edge cases - extreme spacing
		{
			name:  "many consecutive spaces",
			input: "A" + strings.Repeat(" ", 100) + "B",
			want:  "A B",
		},
		{
			name:  "tabs and spaces mixed",
			input: "A\t \t B",
			want:  "A\t \t B", // TrimSpace doesn't replace tabs with spaces
		},
		
		// Edge cases - single character
		{
			name:  "single character",
			input: "A",
			want:  "A",
		},
		{
			name:  "single character with spaces",
			input: "  A  ",
			want:  "A",
		},
		
		// Edge cases - very long strings
		{
			name:  "very long title",
			input: strings.Repeat("A", 1000),
			want:  strings.Repeat("A", 1000),
		},
		{
			name:  "very long with spaces",
			input: "  " + strings.Repeat("A", 1000) + "  ",
			want:  strings.Repeat("A", 1000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeTitle(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeTitle(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestTruncateDescription(t *testing.T) {
	tests := []struct {
		name   string
		desc   string
		maxLen int
		want   string
	}{
		// Happy paths
		{
			name:   "short description no truncate",
			desc:   "Short desc",
			maxLen: 20,
			want:   "Short desc",
		},
		{
			name:   "exact length no truncate",
			desc:   "12345",
			maxLen: 5,
			want:   "12345",
		},
		{
			name:   "truncate with ellipsis",
			desc:   "This is a long description",
			maxLen: 10,
			want:   "This is...",
		},
		{
			name:   "truncate very long text",
			desc:   "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
			maxLen: 20,
			want:   "Lorem ipsum dolor...",
		},
		
		// Edge cases - zero and negative maxLen
		{
			name:   "maxLen zero",
			desc:   "Any text",
			maxLen: 0,
			want:   "",
		},
		{
			name:   "maxLen negative",
			desc:   "Any text",
			maxLen: -1,
			want:   "Any text",
		},
		{
			name:   "maxLen large negative",
			desc:   "Any text",
			maxLen: -1000,
			want:   "Any text",
		},
		
		// Edge cases - very small maxLen
		{
			name:   "maxLen 1",
			desc:   "Hello",
			maxLen: 1,
			want:   "H",
		},
		{
			name:   "maxLen 2",
			desc:   "Hello",
			maxLen: 2,
			want:   "He",
		},
		{
			name:   "maxLen 3",
			desc:   "Hello",
			maxLen: 3,
			want:   "Hel",
		},
		{
			name:   "maxLen 4 adds ellipsis",
			desc:   "Hello World",
			maxLen: 4,
			want:   "H...",
		},
		
		// Edge cases - empty strings
		{
			name:   "empty description",
			desc:   "",
			maxLen: 10,
			want:   "",
		},
		{
			name:   "whitespace only",
			desc:   "   ",
			maxLen: 10,
			want:   "",
		},
		{
			name:   "whitespace truncated",
			desc:   "     ",
			maxLen: 2,
			want:   "",
		},
		
		// Edge cases - whitespace handling
		{
			name:   "leading whitespace removed",
			desc:   "   Hello World",
			maxLen: 8,
			want:   "Hello...",
		},
		{
			name:   "trailing whitespace removed",
			desc:   "Hello World   ",
			maxLen: 20,
			want:   "Hello World",
		},
		
		// Edge cases - unicode and special characters
		{
			name:   "unicode characters",
			desc:   "Hello 世界 World",
			maxLen: 10,
			want:   "Hello \xe4...", // Truncates by byte, not rune
		},
		{
			name:   "emojis",
			desc:   "🔥🔥🔥 Fire sale! 🔥🔥🔥",
			maxLen: 15,
			want:   "🔥🔥🔥...", // Emoji takes multiple bytes
		},
		
		// Edge cases - very large maxLen
		{
			name:   "maxLen larger than desc",
			desc:   "Short",
			maxLen: 1000,
			want:   "Short",
		},
		{
			name:   "maxLen int max",
			desc:   "Test",
			maxLen: 2147483647,
			want:   "Test",
		},
		
		// Edge cases - exactly at ellipsis boundary
		{
			name:   "length equals maxLen exactly",
			desc:   "12345",
			maxLen: 5,
			want:   "12345",
		},
		{
			name:   "length one more than maxLen",
			desc:   "123456",
			maxLen: 5,
			want:   "12...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TruncateDescription(tt.desc, tt.maxLen)
			if got != tt.want {
				t.Errorf("TruncateDescription(%q, %d) = %q, want %q", tt.desc, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestExtractDomain(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr error
	}{
		// Happy paths
		{
			name:    "simple domain",
			url:     "https://example.com/feed",
			want:    "example.com",
			wantErr: nil,
		},
		{
			name:    "subdomain",
			url:     "https://blog.example.com/rss",
			want:    "blog.example.com",
			wantErr: nil,
		},
		{
			name:    "domain with port",
			url:     "https://example.com:8080/feed",
			want:    "example.com:8080",
			wantErr: nil,
		},
		{
			name:    "http scheme",
			url:     "http://example.com",
			want:    "example.com",
			wantErr: nil,
		},
		
		// Edge cases - empty/invalid
		{
			name:    "empty URL",
			url:     "",
			want:    "",
			wantErr: ErrEmptyURL,
		},
		{
			name:    "whitespace only",
			url:     "   ",
			want:    "",
			wantErr: nil, // url.Parse accepts whitespace
		},
		{
			name:    "invalid URL format",
			url:     "not a url",
			want:    "",
			wantErr: nil,
		},
		{
			name:    "malformed URL",
			url:     "://broken",
			want:    "",
			wantErr: ErrInvalidURL,
		},
		
		// Edge cases - special domains
		{
			name:    "localhost",
			url:     "http://localhost/feed",
			want:    "localhost",
			wantErr: nil,
		},
		{
			name:    "IPv4",
			url:     "http://192.168.1.1/feed",
			want:    "192.168.1.1",
			wantErr: nil,
		},
		{
			name:    "IPv6",
			url:     "http://[::1]/feed",
			want:    "[::1]",
			wantErr: nil,
		},
		{
			name:    "IPv6 with port",
			url:     "http://[2001:db8::1]:8080/feed",
			want:    "[2001:db8::1]:8080",
			wantErr: nil,
		},
		
		// Edge cases - URL components
		{
			name:    "with query params",
			url:     "https://example.com/feed?key=value",
			want:    "example.com",
			wantErr: nil,
		},
		{
			name:    "with fragment",
			url:     "https://example.com/feed#section",
			want:    "example.com",
			wantErr: nil,
		},
		{
			name:    "with auth",
			url:     "https://user:pass@example.com/feed",
			want:    "example.com",
			wantErr: nil,
		},
		
		// Edge cases - internationalized domains
		{
			name:    "unicode domain",
			url:     "https://例え.jp/feed",
			want:    "例え.jp",
			wantErr: nil,
		},
		
		// Edge cases - no scheme
		{
			name:    "protocol-relative URL",
			url:     "//example.com/feed",
			want:    "example.com",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractDomain(tt.url)
			if err != tt.wantErr {
				t.Errorf("ExtractDomain(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractDomain(%q) = %q, want %q", tt.url, got, tt.want)
			}
		})
	}
}

func TestIsValidRSSPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		// Happy paths - valid extensions
		{
			name: "rss extension",
			path: "/feed.rss",
			want: true,
		},
		{
			name: "xml extension",
			path: "/feed.xml",
			want: true,
		},
		{
			name: "atom extension",
			path: "/feed.atom",
			want: true,
		},
		
		// Happy paths - valid path names
		{
			name: "rss in path",
			path: "/rss",
			want: true,
		},
		{
			name: "feed in path",
			path: "/feed",
			want: true,
		},
		{
			name: "atom in path",
			path: "/atom",
			want: true,
		},
		{
			name: "blog/rss",
			path: "/blog/rss",
			want: true,
		},
		{
			name: "posts/feed.xml",
			path: "/posts/feed.xml",
			want: true,
		},
		
		// Edge cases - case insensitivity
		{
			name: "RSS uppercase",
			path: "/FEED.RSS",
			want: true,
		},
		{
			name: "XML uppercase",
			path: "/FEED.XML",
			want: true,
		},
		{
			name: "mixed case",
			path: "/FeeD.XmL",
			want: true,
		},
		
		// Edge cases - invalid paths
		{
			name: "empty path",
			path: "",
			want: false,
		},
		{
			name: "just slash",
			path: "/",
			want: false,
		},
		{
			name: "html file",
			path: "/index.html",
			want: false,
		},
		{
			name: "random path",
			path: "/about",
			want: false,
		},
		{
			name: "similar but wrong",
			path: "/feeds",
			want: true, // Contains "feed"
		},
		
		// Edge cases - substring matches
		{
			name: "rss in middle",
			path: "/blog/rss/archive",
			want: true,
		},
		{
			name: "feed in query param",
			path: "/api?type=feed",
			want: false, // Only checks for "/feed", not "feed"
		},
		
		// Edge cases - special characters
		{
			name: "with query params",
			path: "/feed?format=rss",
			want: true,
		},
		{
			name: "with fragment",
			path: "/rss.xml#latest",
			want: true,
		},
		
		// Edge cases - no leading slash
		{
			name: "no leading slash rss",
			path: "feed.rss",
			want: true,
		},
		{
			name: "no leading slash path",
			path: "rss",
			want: false, // Requires "/rss" not "rss"
		},
		
		// Edge cases - whitespace
		{
			name: "whitespace only",
			path: "   ",
			want: false,
		},
		
		// Edge cases - very long paths
		{
			name: "long valid path",
			path: strings.Repeat("/blog", 50) + "/feed.rss",
			want: true,
		},
		{
			name: "long invalid path",
			path: strings.Repeat("/blog", 50) + "/index.html",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidRSSPath(tt.path)
			if got != tt.want {
				t.Errorf("IsValidRSSPath(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr error
	}{
		// Happy paths - already valid
		{
			name:    "https URL unchanged",
			url:     "https://example.com/feed",
			want:    "https://example.com/feed",
			wantErr: nil,
		},
		{
			name:    "http URL unchanged",
			url:     "http://example.com/feed",
			want:    "http://example.com/feed",
			wantErr: nil,
		},
		
		// Happy paths - add https
		{
			name:    "no scheme adds https",
			url:     "example.com/feed",
			want:    "https://example.com/feed",
			wantErr: nil,
		},
		{
			name:    "subdomain no scheme",
			url:     "blog.example.com/rss",
			want:    "https://blog.example.com/rss",
			wantErr: nil,
		},
		{
			name:    "with port no scheme",
			url:     "example.com:8080/feed",
			want:    "https://example.com:8080/feed",
			wantErr: nil,
		},
		
		// Edge cases - whitespace
		{
			name:    "leading whitespace",
			url:     "  https://example.com/feed",
			want:    "https://example.com/feed",
			wantErr: nil,
		},
		{
			name:    "trailing whitespace",
			url:     "https://example.com/feed  ",
			want:    "https://example.com/feed",
			wantErr: nil,
		},
		{
			name:    "whitespace around no-scheme URL",
			url:     "  example.com/feed  ",
			want:    "https://example.com/feed",
			wantErr: nil,
		},
		{
			name:    "empty string",
			url:     "",
			want:    "",
			wantErr: ErrEmptyURL,
		},
		{
			name:    "only whitespace",
			url:     "   ",
			want:    "",
			wantErr: ErrEmptyURL,
		},
		
		// Edge cases - invalid after normalization
		{
			name:    "invalid scheme not changed",
			url:     "ftp://example.com/feed",
			want:    "",
			wantErr: ErrInvalidURL, // Validate should catch this
		},
		{
			name:    "no host after adding scheme",
			url:     "/just/a/path",
			want:    "",
			wantErr: ErrInvalidURL,
		},
		
		// Edge cases - complex URLs
		{
			name:    "with query params",
			url:     "example.com/feed?format=rss&limit=10",
			want:    "https://example.com/feed?format=rss&limit=10",
			wantErr: nil,
		},
		{
			name:    "with fragment",
			url:     "example.com/feed#latest",
			want:    "https://example.com/feed#latest",
			wantErr: nil,
		},
		{
			name:    "with everything",
			url:     "example.com:8080/feed?key=val#section",
			want:    "https://example.com:8080/feed?key=val#section",
			wantErr: nil,
		},
		
		// Edge cases - localhost and IPs
		{
			name:    "localhost no scheme",
			url:     "localhost/feed",
			want:    "https://localhost/feed",
			wantErr: nil,
		},
		{
			name:    "IPv4 no scheme",
			url:     "192.168.1.1/feed",
			want:    "https://192.168.1.1/feed",
			wantErr: nil,
		},
		
		// Edge cases - unicode
		{
			name:    "unicode domain no scheme",
			url:     "例え.jp/feed",
			want:    "https://例え.jp/feed",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeURL(tt.url)
			if err != tt.wantErr {
				t.Errorf("NormalizeURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NormalizeURL(%q) = %q, want %q", tt.url, got, tt.want)
			}
		})
	}
}
