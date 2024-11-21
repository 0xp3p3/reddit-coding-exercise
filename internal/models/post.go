package models

type Post struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Upvotes int    `json:"score"`
}
