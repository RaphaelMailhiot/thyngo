package posts

type PostStore interface {
	ListPosts() []Post
	CreatePost(slug, title, content string) (*Post, error)
	GetPostBySlug(slug string) *Post
	UpdatePostBySlug(slug, title, content string) (*Post, error)
	DeletePostBySlug(slug string) (bool, error)
}
