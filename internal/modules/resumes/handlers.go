package resumes

import (
	"net/http"
	"strconv"

	"thyngo/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (m *ResumesModule) listResumesHandler(c *gin.Context) {
	resumes := m.service.ListResumes()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resumes,
	})
}

func (m *ResumesModule) createResumeHandler(c *gin.Context) {
	var req struct {
		Name       string `json:"name" binding:"required"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		Job        string `json:"job"`
		Github     string `json:"github"`
		Linkedin   string `json:"linkedin"`
		Website    string `json:"website"`
		Visibility string `json:"visibility"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var userID int64
	if id, exists := c.Get(middleware.UserIDKey); exists {
		if uid, ok := id.(int64); ok {
			userID = uid
		}
	}

	resume, err := m.service.CreateResume(userID, req.Name, req.Email, req.Phone, req.Job, req.Github, req.Linkedin, req.Website, req.Visibility)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    resume,
	})
}

func (m *ResumesModule) getResumeHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid resume ID format",
		})
		return
	}

	resume := m.service.GetResume(id)

	if resume == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Resume not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resume,
	})
}

func (m *ResumesModule) updateResumeHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid resume ID format",
		})
		return
	}

	var req struct {
		Name       string `json:"name"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		Job        string `json:"job"`
		Github     string `json:"github"`
		Linkedin   string `json:"linkedin"`
		Website    string `json:"website"`
		Visibility string `json:"visibility"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var userID int64
	if idVal, exists := c.Get(middleware.UserIDKey); exists {
		if uid, ok := idVal.(int64); ok {
			userID = uid
		}
	}

	resume, err := m.service.UpdateResume(userID, id, req.Name, req.Email, req.Phone, req.Job, req.Github, req.Linkedin, req.Website, req.Visibility)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if resume == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Resume not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resume,
	})
}

func (m *ResumesModule) deleteResumeHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid resume ID format",
		})
		return
	}

	var userID int64
	if idVal, exists := c.Get(middleware.UserIDKey); exists {
		if uid, ok := idVal.(int64); ok {
			userID = uid
		}
	}

	deleted, err := m.service.DeleteResume(userID, id)
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
			"error":   "Resume not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Resume deleted.",
	})
}
