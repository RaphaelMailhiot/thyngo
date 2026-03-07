package posts

import (
	"sync"
	"time"
)

type InMemoryStore struct {
	mu        sync.RWMutex
	mem       []Post
	blocks    map[int64]*Block
	contents  map[int64]*Content
	blockID   int64
	contentID int64
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		mem:      []Post{},
		blocks:   make(map[int64]*Block),
		contents: make(map[int64]*Content),
	}
}

func (s *InMemoryStore) ListPosts() []Post {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Post, len(s.mem))
	copy(out, s.mem)
	return out
}

func (s *InMemoryStore) CreatePost(userID *int64, slug, title, visibility string) (*Post, error) {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()

	if visibility == "" {
		visibility = "public"
	}

	id := int64(len(s.mem) + 1)
	p := Post{
		ID:         id,
		UserID:     userID,
		Slug:       slug,
		Title:      title,
		Visibility: visibility,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	s.mem = append(s.mem, p)
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

func (s *InMemoryStore) GetPostBySlugWithContents(slug string) *Post {
	p := s.GetPostBySlug(slug)
	if p == nil {
		return nil
	}
	p.Contents = s.ListContentsByPostID(p.ID)
	return p
}

func (s *InMemoryStore) UpdatePostBySlug(slug, title, visibility string) (*Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if visibility == "" {
		visibility = "public"
	}

	for i := range s.mem {
		if s.mem[i].Slug == slug {
			s.mem[i].Title = title
			s.mem[i].Visibility = visibility
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

// Block management

func (s *InMemoryStore) CreateBlock(blockType string, data BlockData) (*Block, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockID++
	now := time.Now()
	b := &Block{
		ID:        s.blockID,
		Type:      blockType,
		CreatedAt: now,
		UpdatedAt: now,
		Data:      data,
	}
	s.blocks[b.ID] = b
	return b, nil
}

func (s *InMemoryStore) GetBlock(blockID int64) *Block {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if b, ok := s.blocks[blockID]; ok {
		return b
	}
	return nil
}

func (s *InMemoryStore) UpdateBlock(blockID int64, data BlockData) (*Block, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if b, ok := s.blocks[blockID]; ok {
		b.UpdatedAt = time.Now()
		b.Data = data
		return b, nil
	}
	return nil, nil
}

func (s *InMemoryStore) DeleteBlock(blockID int64) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.blocks[blockID]; ok {
		delete(s.blocks, blockID)
		return true, nil
	}
	return false, nil
}

// Content management

func (s *InMemoryStore) CreateContent(parentID int64, order int, blockType string, blockID int64) (*Content, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.contentID++
	now := time.Now()
	c := &Content{
		ID:        s.contentID,
		ParentID:  parentID,
		Order:     order,
		Type:      blockType,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if b, ok := s.blocks[blockID]; ok {
		c.Block = b
	}
	s.contents[c.ID] = c
	return c, nil
}

func (s *InMemoryStore) ListContentsByPostID(postID int64) []Content {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []Content
	for _, c := range s.contents {
		if c.ParentID == postID {
			out = append(out, *c)
		}
	}
	return out
}

func (s *InMemoryStore) UpdateContent(contentID int64, order int) (*Content, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c, ok := s.contents[contentID]; ok {
		c.Order = order
		c.UpdatedAt = time.Now()
		return c, nil
	}
	return nil, nil
}

func (s *InMemoryStore) DeleteContent(contentID int64) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.contents[contentID]; ok {
		delete(s.contents, contentID)
		return true, nil
	}
	return false, nil
}
