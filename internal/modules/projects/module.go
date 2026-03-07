package projects

import (
	"thyngo/internal/middleware"

	"github.com/gin-gonic/gin"
)

type ProjectsModule struct {
	service *Service
}

func New() *ProjectsModule {
	return &ProjectsModule{
		service: NewService(),
	}
}

func (m *ProjectsModule) Name() string {
	return "projects"
}

func (m *ProjectsModule) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("", m.listProjectsHandler)
	router.GET("/:slug", m.getProjectHandler)

	protected := router.Group("")
	protected.Use(middleware.RequireAuth())
	{
		protected.POST("", m.createProjectsHandler)
		protected.PUT("/:slug", m.updateProjectHandler)
		protected.DELETE("/:slug", m.deleteProjectHandler)
	}
}
