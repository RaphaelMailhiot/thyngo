package media

import (
	"net/http"
	"strconv"

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

// Movies Handlers

func (m *MediaModule) listMoviesHandler(c *gin.Context) {
	movies := m.service.ListMovies()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    movies,
	})
}

func (m *MediaModule) createMovieHandler(c *gin.Context) {
	var movie Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if id, exists := c.Get(middleware.UserIDKey); exists {
		if uid, ok := id.(int64); ok {
			movie.UserID = &uid
		}
	}

	if err := m.service.CreateMovie(&movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": movie})
}

func (m *MediaModule) updateMovieHandler(c *gin.Context) {
	idStr := c.Param("id")
	var req Movie
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	id, _ := strconv.ParseInt(idStr, 10, 64)
	req.ID = id

	if err := m.service.UpdateMovie(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": req})
}

// Series Handlers

func (m *MediaModule) listSeriesHandler(c *gin.Context) {
	series := m.service.ListSeries()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    series,
	})
}

func (m *MediaModule) createSeriesHandler(c *gin.Context) {
	var sr Series
	if err := c.ShouldBindJSON(&sr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if id, exists := c.Get(middleware.UserIDKey); exists {
		if uid, ok := id.(int64); ok {
			sr.UserID = &uid
		}
	}

	if err := m.service.CreateSeries(&sr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": sr})
}

func (m *MediaModule) updateSeriesHandler(c *gin.Context) {
	idStr := c.Param("id")
	var req Series
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	id, _ := strconv.ParseInt(idStr, 10, 64)
	req.ID = id

	if err := m.service.UpdateSeries(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": req})
}

// Albums Handlers

func (m *MediaModule) listAlbumsHandler(c *gin.Context) {
	albums := m.service.ListAlbums()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    albums,
	})
}

func (m *MediaModule) createAlbumHandler(c *gin.Context) {
	var a MusicAlbum
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if id, exists := c.Get(middleware.UserIDKey); exists {
		if uid, ok := id.(int64); ok {
			a.UserID = &uid
		}
	}

	if err := m.service.CreateAlbum(&a); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": a})
}

func (m *MediaModule) updateAlbumHandler(c *gin.Context) {
	idStr := c.Param("id")
	var req MusicAlbum
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	id, _ := strconv.ParseInt(idStr, 10, 64)
	req.ID = id

	if err := m.service.UpdateAlbum(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": req})
}

// External Metadata Handlers

func (m *MediaModule) searchExternalMoviesHandler(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Query is required"})
		return
	}

	client := NewTMDBClient()
	results, err := client.SearchMovies(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

func (m *MediaModule) searchExternalSeriesHandler(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Query is required"})
		return
	}

	client := NewTMDBClient()
	results, err := client.SearchSeries(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

func (m *MediaModule) searchExternalAlbumsHandler(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Query is required"})
		return
	}

	mbClient := NewMusicBrainzClient()
	results, err := mbClient.SearchAlbums(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

func (m *MediaModule) getExternalAlbumArtworkHandler(c *gin.Context) {
	mbid := c.Param("mbid")
	if mbid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "MBID is required"})
		return
	}

	fClient := NewFanartClient()
	urls, err := fClient.GetAlbumArtwork(mbid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": urls})
}

func (m *MediaModule) getExternalMovieArtworkHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID is required"})
		return
	}

	fClient := NewFanartClient()
	urls, err := fClient.GetMovieArtwork(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": urls})
}

func (m *MediaModule) getExternalSeriesArtworkHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID is required"})
		return
	}

	fClient := NewFanartClient()
	urls, err := fClient.GetSeriesArtwork(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": urls})
}

func (m *MediaModule) importExternalMediaHandler(c *gin.Context) {
	var req struct {
		Title      string `json:"title"`
		URL        string `json:"url"`
		Visibility string `json:"visibility"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if req.URL == "" || req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "URL and Title are required"})
		return
	}

	var userID *int64
	if id, exists := c.Get(middleware.UserIDKey); exists {
		if uid, ok := id.(int64); ok {
			userID = &uid
		}
	}

	media, err := m.service.ImportMediaFromURL(userID, req.Title, req.URL, req.Visibility)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": media})
}
