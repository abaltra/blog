package post

import "testing"

func TestSlugCreation(t *testing.T) {
	p := Post{
		Title: "this has spaces",
	}

	expected := "this-has-spaces"

	p.BuildSlug()

	if expected != p.Slug {
		t.Errorf("Expected %s to equal %s", p.Slug, expected)
	}
}
