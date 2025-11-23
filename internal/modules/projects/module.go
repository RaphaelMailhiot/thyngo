package projects

import "github.com/gin-gonic/gin"

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
	//router.GET("", m.listProjectsHandler)
	//router.POST("", m.createProjectsHandler)
	//router.GET("/:slug", m.getProjectHandler)
	//router.PUT("/:slug", m.updateProjectHandler)
	//router.DELETE("/:slug", m.deleteProjectHandler)
}
