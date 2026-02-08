package rss

import (
	"errors"
	"net/url"
	"strings"
)

var (
	ErrEmptyURL       = errors.New("URL cannot be empty")
	ErrInvalidURL     = errors.New("invalid URL format")
	ErrInvalidScheme  = errors.New("URL must use http or https scheme")
	ErrEmptyTitle     = errors.New("title cannot be empty")
	ErrTitleTooLong   = errors.New("title exceeds maximum length")
	ErrInvalidXML     = errors.New("invalid XML content")
)

const MaxTitleLength = 500

// ValidateURL checks if a URL is valid for RSS feed fetching
func ValidateURL(rawURL string) error {
	if rawURL == "" {
		return ErrEmptyURL
	}

	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ErrEmptyURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ErrInvalidURL
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return ErrInvalidScheme
	}

	if parsedURL.Host == "" {
		return ErrInvalidURL
	}

	return nil
}

// ValidateTitle checks if a feed title is valid
func ValidateTitle(title string) error {
	title = strings.TrimSpace(title)
	
	if title == "" {
		return ErrEmptyTitle
	}

	if len(title) > MaxTitleLength {
		return ErrTitleTooLong
	}

	return nil
}

// SanitizeTitle cleans and normalizes a feed title
func SanitizeTitle(title string) string {
	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, "\n", " ")
	title = strings.ReplaceAll(title, "\r", " ")
	title = strings.ReplaceAll(title, "\t", " ")
	
	// Replace multiple spaces with single space
	for strings.Contains(title, "  ") {
		title = strings.ReplaceAll(title, "  ", " ")
	}
	
	if len(title) > MaxTitleLength {
		title = title[:MaxTitleLength]
	}
	
	return title
}

// NormalizeURL normalizes a URL by removing fragments and ensuring scheme
func NormalizeURL(rawURL string) (string, error) {
	if err := ValidateURL(rawURL); err != nil {
		return "", err
	}

	parsedURL, _ := url.Parse(strings.TrimSpace(rawURL))
	
	// Remove fragment
	parsedURL.Fragment = ""
	
	// Ensure lowercase scheme and host
	parsedURL.Scheme = strings.ToLower(parsedURL.Scheme)
	parsedURL.Host = strings.ToLower(parsedURL.Host)
	
	return parsedURL.String(), nil
}

// IsValidXMLStart checks if content starts with valid XML declaration
func IsValidXMLStart(content string) bool {
	content = strings.TrimSpace(content)
	if len(content) == 0 {
		return false
	}
	
	return strings.HasPrefix(content, "<?xml") || 
	       strings.HasPrefix(content, "<rss") || 
	       strings.HasPrefix(content, "<feed")
}
