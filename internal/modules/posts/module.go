package posts

import (
	"thyngo/internal/middleware"

	"github.com/gin-gonic/gin"
)

type PostsModule struct {
	service PostStore
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
	// Post routes
	router.GET("", m.listPostsHandler)
	router.GET("/:slug", m.getPostHandler)

	protected := router.Group("")
	protected.Use(middleware.RequireAuth())
	{
		protected.POST("", m.createPostsHandler)
		protected.PUT("/:slug", m.updatePostHandler)
		protected.DELETE("/:slug", m.deletePostHandler)

		// Content routes for posts
		protected.POST("/:slug/contents", m.createContentHandler)
		protected.PUT("/contents/:contentID", m.updateContentHandler)
		protected.DELETE("/contents/:contentID", m.deleteContentHandler)

		// Block routes
		protected.POST("/blocks", m.createBlockHandler)
		protected.PUT("/blocks/:blockID", m.updateBlockHandler)
		protected.DELETE("/blocks/:blockID", m.deleteBlockHandler)
	}

	router.GET("/:slug/contents", m.listContentsHandler)
	router.GET("/blocks/:blockID", m.getBlockHandler)
}
