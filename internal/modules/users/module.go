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
	router.POST("/register", m.registerHandler)
	router.POST("/login", m.loginHandler)
}
