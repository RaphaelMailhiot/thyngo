package posts

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// createBlockHandler creates a new block
func (m *PostsModule) createBlockHandler(c *gin.Context) {
	var req struct {
		Type string      `json:"type" binding:"required"`
		Data interface{} `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	block, err := m.service.CreateBlock(req.Type, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    block,
	})
}

// getBlockHandler retrieves a block by ID
func (m *PostsModule) getBlockHandler(c *gin.Context) {
	blockID, err := strconv.ParseInt(c.Param("blockID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid block ID",
		})
		return
	}

	block := m.service.GetBlock(blockID)
	if block == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Block not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    block,
	})
}

// updateBlockHandler updates a block
func (m *PostsModule) updateBlockHandler(c *gin.Context) {
	blockID, err := strconv.ParseInt(c.Param("blockID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid block ID",
		})
		return
	}

	var req struct {
		Data interface{} `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	block, err := m.service.UpdateBlock(blockID, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if block == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Block not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    block,
	})
}

// deleteBlockHandler deletes a block
func (m *PostsModule) deleteBlockHandler(c *gin.Context) {
	blockID, err := strconv.ParseInt(c.Param("blockID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid block ID",
		})
		return
	}

	deleted, err := m.service.DeleteBlock(blockID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Block not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Block deleted",
	})
}

// createContentHandler creates content for a post
func (m *PostsModule) createContentHandler(c *gin.Context) {
	slug := c.Param("slug")
	post := m.service.GetPostBySlug(slug)
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Post not found",
		})
		return
	}

	var req struct {
		Order   *int   `json:"order" binding:"required"`
		Type    string `json:"type" binding:"required"`
		BlockID int64  `json:"block_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	content, err := m.service.CreateContent(post.ID, *req.Order, req.Type, req.BlockID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    content,
	})
}

// listContentsHandler retrieves all contents for a post
func (m *PostsModule) listContentsHandler(c *gin.Context) {
	slug := c.Param("slug")
	post := m.service.GetPostBySlug(slug)
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Post not found",
		})
		return
	}

	contents := m.service.ListContentsByPostID(post.ID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    contents,
	})
}

// updateContentHandler updates content order
func (m *PostsModule) updateContentHandler(c *gin.Context) {
	contentID, err := strconv.ParseInt(c.Param("contentID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid content ID",
		})
		return
	}

	var req struct {
		Order *int `json:"order" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	content, err := m.service.UpdateContent(contentID, *req.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if content == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Content not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    content,
	})
}

// deleteContentHandler deletes content
func (m *PostsModule) deleteContentHandler(c *gin.Context) {
	contentID, err := strconv.ParseInt(c.Param("contentID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid content ID",
		})
		return
	}

	deleted, err := m.service.DeleteContent(contentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Content not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Content deleted",
	})
}
