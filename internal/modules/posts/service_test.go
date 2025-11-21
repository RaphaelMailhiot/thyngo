package posts

import "testing"

func TestListPosts(t *testing.T) {
	s := NewService()
	posts := s.ListPosts()
	if len(posts) < 1 {
		t.Fatalf("expected at least one post, got %d", len(posts))
	}
}

func TestGetPostBySlug(t *testing.T) {
	s := NewService()
	p := s.GetPostBySlug("first-post")
	if p == nil || p.Slug != "first-post" {
		t.Fatalf("expected to find first-post")
	}
	p2 := s.GetPostBySlug("not-found")
	if p2 != nil {
		t.Fatalf("expected nil for unknown slug")
	}
}
