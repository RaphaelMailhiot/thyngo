package posts

type Post struct {
	ID      string `json:"id"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Service struct {
	// repository PostRepository
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ListPosts() []Post {
	// Placeholder implementation
	posts := []Post{
		{ID: "1", Slug: "first-post", Title: "First Post", Content: "This is the content of the first post."},
		{ID: "2", Slug: "second-post", Title: "Second Post", Content: "This is the content of the second post."},
	}
	return posts
}

func (s *Service) GetPostBySlug(slug string) *Post {
	// Placeholder implementation
	if slug == "first-post" {
		return &Post{ID: "1", Slug: "first-post", Title: "First Post", Content: "This is the content of the first post."}
	} else if slug == "second-post" {
		return &Post{ID: "2", Slug: "second-post", Title: "Second Post", Content: "This is the content of the second post."}
	}
	return nil
}
