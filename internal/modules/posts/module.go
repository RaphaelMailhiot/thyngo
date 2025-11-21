package posts

import "github.com/gin-gonic/gin"

type PostsModule struct {
	service *Service
}

func New() *PostsModule {
	return &PostsModule{
		service: NewService(),
	}
}

func (m *PostsModule) Name() string {
	return "posts"
}

func (m *PostsModule) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("", m.listPostsHandler)
	router.GET("/:slug", m.getPostHandler)
}
