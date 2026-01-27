package models

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


