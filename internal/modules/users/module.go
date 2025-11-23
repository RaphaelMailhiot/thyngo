package users

import "github.com/gin-gonic/gin"

type UsersModule struct {
	service *Service
}

func New() *UsersModule {
	return &UsersModule{
		service: NewService(),
	}
}

func (m *UsersModule) Name() string {
	return "users"
}

func (m *UsersModule) RegisterRoutes(router *gin.RouterGroup) {
	//router.GET("", m.listProjectsHandler)
	//router.POST("", m.createProjectsHandler)
	//router.GET("/:slug", m.getProjectHandler)
	//router.PUT("/:slug", m.updateProjectHandler)
	//router.DELETE("/:slug", m.deleteProjectHandler)
}
