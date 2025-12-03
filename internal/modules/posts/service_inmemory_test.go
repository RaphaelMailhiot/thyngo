package posts

import (
	"sync"
	"time"
)

type InMemoryStore struct {
	mu  sync.RWMutex
	mem []Post
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{mem: []Post{}}
}

func (s *InMemoryStore) ListPosts() []Post {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Post, len(s.mem))
	copy(out, s.mem)
	return out
}

func (s *InMemoryStore) CreatePost(slug, title, content string) (*Post, error) {
	now := time.Now()
	p := Post{
		Slug:      slug,
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.mu.Lock()
	s.mem = append(s.mem, p)
	s.mu.Unlock()
	return &p, nil
}

func (s *InMemoryStore) GetPostBySlug(slug string) *Post {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.mem {
		if p.Slug == slug {
			c := p
			return &c
		}
	}
	return nil
}

func (s *InMemoryStore) UpdatePostBySlug(slug, title, content string) (*Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.mem {
		if s.mem[i].Slug == slug {
			s.mem[i].Title = title
			s.mem[i].UpdatedAt = time.Now()
			c := s.mem[i]
			return &c, nil
		}
	}
	return nil, nil
}

func (s *InMemoryStore) DeletePostBySlug(slug string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.mem {
		if s.mem[i].Slug == slug {
			s.mem = append(s.mem[:i], s.mem[i+1:]...)
			return true, nil
		}
	}
	return false, nil
}
