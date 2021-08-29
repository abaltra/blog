package post

import (
	"strings"
	"time"
)

type Post struct {
	ID          int       `json:"ID"`
	Slug        string    `json:"Slug"`
	Title       string    `json:"Title"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
	PublishedAt time.Time `json:"PublishedAt"`
	Version     int       `json:"Version"`
	AuthorID    string    `json:"AuthorID"`
	Abstract    string    `json:"Abstract"`
	ContentRaw  string    `json:"ContentRaw"`
	IsPublished bool      `json:"IsPublished"`
}

func (p *Post) BuildSlug() {
	s := strings.ToLower(p.Title)
	s = strings.ReplaceAll(s, " ", "-")
	p.Slug = s
}
