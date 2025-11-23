package posts

import "testing"

func TestListPosts(t *testing.T) {
	s := NewInMemoryStore()
	_, _ = s.CreatePost("first-post", "First", "content")
	posts := s.ListPosts()
	if len(posts) < 1 {
		t.Fatalf("expected at least one post, got %d", len(posts))
	}
}

func TestGetCreateUpdateDeletePost(t *testing.T) {
	s := NewInMemoryStore()

	// Create
	created, err := s.CreatePost("first-post", "First", "content")
	if err != nil || created == nil {
		t.Fatalf("create failed: %v", err)
	}
	if created.Slug != "first-post" {
		t.Fatalf("expected slug first-post, got %s", created.Slug)
	}

	// Get existing
	p := s.GetPostBySlug("first-post")
	if p == nil || p.Slug != "first-post" {
		t.Fatalf("expected to find first-post")
	}

	// Get not found
	p2 := s.GetPostBySlug("not-found")
	if p2 != nil {
		t.Fatalf("expected nil for unknown slug")
	}

	// Update
	updated, err := s.UpdatePostBySlug("first-post", "Updated", "new content")
	if err != nil || updated == nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Title != "Updated" {
		t.Fatalf("expected title Updated, got %s", updated.Title)
	}

	// Delete
	deleted, err := s.DeletePostBySlug("first-post")
	if err != nil || !deleted {
		t.Fatalf("delete failed: %v", err)
	}
	after := s.GetPostBySlug("first-post")
	if after != nil {
		t.Fatalf("expected post to be deleted")
	}
}
