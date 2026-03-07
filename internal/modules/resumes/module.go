package resumes

import (
	"thyngo/internal/middleware"

	"github.com/gin-gonic/gin"
)

type ResumesModule struct {
	service *Service
}

func New() *ResumesModule {
	return &ResumesModule{
		service: NewService(),
	}
}

func (m *ResumesModule) Name() string {
	return "resumes"
}

func (m *ResumesModule) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("", m.listResumesHandler)
	router.GET("/:id", m.getResumeHandler)

	protected := router.Group("")
	protected.Use(middleware.RequireAuth())
	{
		protected.POST("", m.createResumeHandler)
		protected.PUT("/:id", m.updateResumeHandler)
		protected.DELETE("/:id", m.deleteResumeHandler)
	}
}
