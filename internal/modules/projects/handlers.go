package projects

import (
	"net/http"

	"thyngo/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (m *ProjectsModule) listProjectsHandler(c *gin.Context) {
	projects := m.service.ListProjects()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    projects,
	})
}

func (m *ProjectsModule) createProjectsHandler(c *gin.Context) {
	var req struct {
		Slug       string `json:"slug" binding:"required"`
		Title      string `json:"title" binding:"required"`
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

	project, err := m.service.CreateProject(userID, req.Slug, req.Title, req.Visibility)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    project,
	})
}

func (m *ProjectsModule) getProjectHandler(c *gin.Context) {
	slug := c.Param("slug")

	project := m.service.GetProjectBySlug(slug)

	if project == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Project not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    project,
	})
}

func (m *ProjectsModule) updateProjectHandler(c *gin.Context) {
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

	project, err := m.service.UpdateProjectBySlug(userID, slug, req.Title, req.Visibility)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if project == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Project not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    project,
	})
}

func (m *ProjectsModule) deleteProjectHandler(c *gin.Context) {
	slug := c.Param("slug")
	var userID *int64
	if id, exists := c.Get(middleware.UserIDKey); exists {
		if uid, ok := id.(int64); ok {
			userID = &uid
		}
	}

	deleted, err := m.service.DeleteProjectBySlug(userID, slug)
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
			"error":   "Project not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Project with slug " + slug + " deleted.",
	})
}
