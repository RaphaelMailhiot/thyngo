package resume

import "github.com/gin-gonic/gin"

type ResumeModule struct {
	service *Service
}

func New() *ResumeModule {
	return &ResumeModule{
		service: NewService(),
	}
}

func (m *ResumeModule) Name() string {
	return "resume"
}

func (m *ResumeModule) RegisterRoutes(router *gin.RouterGroup) {
	//router.GET("", m.listResumeHandler)
	//router.POST("", m.createResumeHandler)
	//router.PUT("", m.updateResumeHandler)
	//router.DELETE("", m.deleteResumeHandler)
}
