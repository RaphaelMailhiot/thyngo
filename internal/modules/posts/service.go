package posts

import (
	"time"
)

// Block types
const (
	BlockTypeImage  = "image"
	BlockTypeTitle  = "title"
	BlockTypeText   = "text"
	BlockTypeCode   = "code"
	BlockTypeCustom = "custom"
)

type Block struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"` // image, title, text, code, custom
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Data      BlockData `json:"data,omitempty"`
}

type BlockData interface{}

type ImageBlockData struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type TitleBlockData struct {
	Level   int    `json:"level"`
	Content string `json:"content"`
}

type TextBlockData struct {
	Content string `json:"content"`
}

type CodeBlockData struct {
	Language string `json:"language"`
	Content  string `json:"content"`
}

type CustomBlockData struct {
	Payload map[string]interface{} `json:"payload"`
}

type Content struct {
	ID        int64     `json:"id"`
	ParentID  int64     `json:"parent_id"`
	Order     int       `json:"order"`
	Type      string    `json:"type"`
	Block     *Block    `json:"block,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Post struct {
	ID         int64     `json:"id"`
	UserID     *int64    `json:"user_id,omitempty"`
	Slug       string    `json:"slug"`
	Title      string    `json:"title"`
	Visibility string    `json:"visibility"`
	Contents   []Content `json:"contents,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

func NewService() PostStore {
	return NewPostgresStore()
}
