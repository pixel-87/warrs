package models

import "strings"

type Post struct {
	Title   string
	Content string
	Link    string
	ID      int
	Read    bool
}

// An entire Feed
type Feed struct {
	Title string
	URL   string
	Posts []Post
	ID    int
}

// HasUnreadPosts returns true if the feed has any unread posts
func (f *Feed) HasUnreadPosts() bool {
	for _, post := range f.Posts {
		if !post.Read {
			return true
		}
	}
	return false
}

// UnreadCount returns the number of unread posts in the feed
func (f *Feed) UnreadCount() int {
	count := 0
	for _, post := range f.Posts {
		if !post.Read {
			count++
		}
	}
	return count
}

// IsValid checks if a Post has the minimum required fields
func (p *Post) IsValid() bool {
	return strings.TrimSpace(p.Title) != "" && strings.TrimSpace(p.Link) != ""
}

// Sanitize returns a copy of the Post with trimmed whitespace
func (p *Post) Sanitize() Post {
	return Post{
		ID:      p.ID,
		Title:   strings.TrimSpace(p.Title),
		Content: strings.TrimSpace(p.Content),
		Link:    strings.TrimSpace(p.Link),
		Read:    p.Read,
	}
}

