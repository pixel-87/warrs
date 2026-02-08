package rss

import (
	"errors"
)

var (
	ErrInvalidURL  = errors.New("invalid URL")
	ErrEmptyURL    = errors.New("URL cannot be empty")
	ErrFetchFailed = errors.New("failed to fetch feed")
)
