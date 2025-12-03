package resumes

import "github.com/gin-gonic/gin"

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
	//router.GET("", m.listResumeHandler)
	//router.POST("", m.createResumeHandler)
	//router.PUT("", m.updateResumeHandler)
	//router.DELETE("", m.deleteResumeHandler)
}
