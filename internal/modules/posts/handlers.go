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

func (m *PostsModule) createPostsHandler(c *gin.Context) {
	var req struct {
		Slug    string `json:"slug" binding:"required"`
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	post, err := m.service.CreatePost(req.Slug, req.Title, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    post,
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

func (m *PostsModule) deletePostHandler(c *gin.Context) {
	slug := c.Param("slug")
	post := m.service.DeletePostBySlug(slug)
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Post not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post with slug " + slug + " deleted (not really, this is a placeholder).",
	})
}
