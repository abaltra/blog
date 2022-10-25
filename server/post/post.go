package post

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID           string    `json:"ID"`
	Slug         string    `json:"Slug"`
	Title        string    `json:"Title"`
	CreatedAt    time.Time `json:"CreatedAt"`
	UpdatedAt    time.Time `json:"UpdatedAt"`
	PublishedAt  time.Time `json:"PublishedAt"`
	Version      int       `json:"Version"`
	AuthorID     string    `json:"AuthorID"`
	Abstract     string    `json:"Abstract"`
	ContentRaw   string    `json:"ContentRaw"`
	IsPublished  bool      `json:"IsPublished"`
	LastEditedBy string    `json:"LastEditedBy"`
}

type CreatePostRequest struct {
	Title      string `json:"Title"`
	Abstract   string `json:"Abstract"`
	ContentRaw string `json:"ContentRaw"`
}

func NewPost(author string, pr CreatePostRequest) *Post {
	return &Post{
		CreatedAt:  time.Now(),
		ID:         uuid.New().String(),
		Version:    1,
		AuthorID:   author,
		Title:      pr.Title,
		Abstract:   pr.Abstract,
		ContentRaw: pr.ContentRaw,
		Slug:       BuildSlug(pr.Title),
	}
}

func BuildSlug(title string) string {
	s := strings.ToLower(title)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
