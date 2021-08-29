package post

import (
	"database/sql"
	"fmt"

	"github.com/abaltra/blog/server/config"
)

type Repository struct {
	DB     *sql.DB
	Config *config.Config
}

func (m *Repository) Create() *Post {
	fmt.Println("Creating a post")
	return &Post{
		Slug: "oh, a new one!",
	}
}

func (m *Repository) Update(postID int, newContent string, publish bool) *Post {
	fmt.Printf("Updating post %d with content %s. Should we publish? %t\n", postID, newContent, publish)
	return &Post{
		ContentRaw:  newContent,
		IsPublished: publish,
	}
}

func (m *Repository) DeleteBySlug(slug string) {
	fmt.Printf("Deleting post by slug: %s\n", slug)
}

func (m *Repository) DeleteByID(id string) {
	fmt.Printf("Deleting post by ID: %s\n", id)
}

func (m *Repository) List(from int, size int) []*Post {
	fmt.Printf("Listing posts. From %d, page size %d\n", from, size)
	return []*Post{
		{
			Slug: "here's one",
		},
		{
			Slug: "here's another one",
		},
	}
}

func (m *Repository) GetBySlug(slug string) *Post {
	fmt.Printf("Getting post by slug %s\n", slug)
	return &Post{
		Slug: "gotten by slug!",
	}
}
