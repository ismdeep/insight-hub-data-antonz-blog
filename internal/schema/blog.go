package schema

import "time"

// Blog struct
type Blog struct {
	Source  string
	Title   string
	Link    string
	Content string
	Date    time.Time
	Author  string
}
