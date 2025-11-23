package media

import "github.com/gin-gonic/gin"

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
	//router.GET("", m.listProjectsHandler)
	//router.POST("", m.createProjectsHandler)
	//router.GET("/:slug", m.getProjectHandler)
	//router.PUT("/:slug", m.updateProjectHandler)
	//router.DELETE("/:slug", m.deleteProjectHandler)
}
