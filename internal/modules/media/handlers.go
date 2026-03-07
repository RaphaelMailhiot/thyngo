package media

import (
	"net/http"

	"thyngo/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (m *MediaModule) listMediaHandler(c *gin.Context) {
	mediaList := m.service.ListMedia()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    mediaList,
	})
}

func (m *MediaModule) createMediaHandler(c *gin.Context) {
	// Multipart form
	form, _ := c.MultipartForm()
	files := form.File["file"]

	slug := c.PostForm("slug")
	title := c.PostForm("title")
	mediaType := c.PostForm("type")
	visibility := c.PostForm("visibility")

	if slug == "" || title == "" || len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Missing required fields or file",
		})
		return
	}

	file := files[0]
	// Save the file
	filePath := "uploads/" + slug + "_" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to save file: " + err.Error(),
		})
		return
	}

	var userID *int64
	if id, exists := c.Get(middleware.UserIDKey); exists {
		if uid, ok := id.(int64); ok {
			userID = &uid
		}
	}

	// In a real app, URL would be the public URL
	url := "/" + filePath

	media, err := m.service.CreateMedia(userID, slug, title, mediaType, url, file.Filename, file.Size, file.Header.Get("Content-Type"), visibility)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    media,
	})
}

func (m *MediaModule) getMediaHandler(c *gin.Context) {
	slug := c.Param("slug")

	media := m.service.GetMediaBySlug(slug)

	if media == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Media not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    media,
	})
}

func (m *MediaModule) updateMediaHandler(c *gin.Context) {
	slug := c.Param("slug")
	var req struct {
		Title      string `json:"title"`
		Visibility string `json:"visibility"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var userID *int64
	if id, exists := c.Get(middleware.UserIDKey); exists {
		if uid, ok := id.(int64); ok {
			userID = &uid
		}
	}

	media, err := m.service.UpdateMediaBySlug(userID, slug, req.Title, req.Visibility)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if media == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Media not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    media,
	})
}

func (m *MediaModule) deleteMediaHandler(c *gin.Context) {
	slug := c.Param("slug")
	var userID *int64
	if id, exists := c.Get(middleware.UserIDKey); exists {
		if uid, ok := id.(int64); ok {
			userID = &uid
		}
	}

	deleted, err := m.service.DeleteMediaBySlug(userID, slug)
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
			"error":   "Media not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Media with slug " + slug + " deleted.",
	})
}
