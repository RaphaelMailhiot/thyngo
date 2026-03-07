package media

import (
	"thyngo/internal/middleware"

	"github.com/gin-gonic/gin"
)

type MediaModule struct {
	service *Service
}

func New() *MediaModule {
	return &MediaModule{
		service: NewService(),
	}
}

func (m *MediaModule) Name() string {
	return "media"
}

func (m *MediaModule) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("", m.listMediaHandler)
	router.GET("/:slug", m.getMediaHandler)

	protected := router.Group("")
	protected.Use(middleware.RequireAuth())
	{
		protected.POST("", m.createMediaHandler)
		protected.PUT("/:slug", m.updateMediaHandler)
		protected.DELETE("/:slug", m.deleteMediaHandler)
	}
}
