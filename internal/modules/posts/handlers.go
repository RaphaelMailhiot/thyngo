package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *PostsModule) listPostsHandler(c *gin.Context) {
	posts := m.service.ListPosts()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    posts,
	})
}

func (m *PostsModule) getPostHandler(c *gin.Context) {
	slug := c.Param("slug")
	post := m.service.GetPostBySlug(slug)
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Post not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    post,
	})
}
