package posts

import (
	"time"
)

type Post struct {
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func NewService() PostStore {
	return NewPostgresStore()
}
