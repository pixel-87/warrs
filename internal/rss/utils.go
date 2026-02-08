package rss

import (
	"net/url"
	"strings"
)

// ValidateURL checks if a URL is valid for RSS feed fetching
// This is a pure function that only validates the URL structure
func ValidateURL(urlStr string) error {
	if urlStr == "" {
		return ErrEmptyURL
	}

	urlStr = strings.TrimSpace(urlStr)
	if urlStr == "" {
		return ErrEmptyURL
	}

	parsed, err := url.Parse(urlStr)
	if err != nil {
		return ErrInvalidURL
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ErrInvalidURL
	}

	if parsed.Host == "" {
		return ErrInvalidURL
	}

	return nil
}

// SanitizeTitle removes leading/trailing whitespace and normalizes internal spacing
func SanitizeTitle(title string) string {
	title = strings.TrimSpace(title)
	// Replace multiple spaces with single space
	for strings.Contains(title, "  ") {
		title = strings.Replace(title, "  ", " ", -1)
	}
	return title
}

// TruncateDescription truncates a description to maxLen characters, adding ellipsis if needed
func TruncateDescription(desc string, maxLen int) string {
	if maxLen < 0 {
		return desc
	}

	if maxLen == 0 {
		return ""
	}

	desc = strings.TrimSpace(desc)
	if len(desc) <= maxLen {
		return desc
	}

	if maxLen <= 3 {
		return desc[:maxLen]
	}

	return desc[:maxLen-3] + "..."
}

// ExtractDomain extracts the domain from a URL string
func ExtractDomain(urlStr string) (string, error) {
	if urlStr == "" {
		return "", ErrEmptyURL
	}

	parsed, err := url.Parse(urlStr)
	if err != nil {
		return "", ErrInvalidURL
	}

	return parsed.Host, nil
}

// IsValidRSSPath checks if a path looks like it could be an RSS feed
func IsValidRSSPath(path string) bool {
	if path == "" {
		return false
	}

	path = strings.ToLower(path)
	validExtensions := []string{".rss", ".xml", ".atom"}
	validPaths := []string{"/rss", "/feed", "/atom"}

	for _, ext := range validExtensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	for _, validPath := range validPaths {
		if strings.Contains(path, validPath) {
			return true
		}
	}

	return false
}

// NormalizeURL adds https:// scheme if missing and validates
func NormalizeURL(urlStr string) (string, error) {
	urlStr = strings.TrimSpace(urlStr)
	
	if urlStr == "" {
		return "", ErrEmptyURL
	}

	// If no scheme, add https://
	if !strings.Contains(urlStr, "://") {
		urlStr = "https://" + urlStr
	}

	// Validate the normalized URL
	if err := ValidateURL(urlStr); err != nil {
		return "", err
	}

	return urlStr, nil
}
