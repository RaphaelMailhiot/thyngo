package server

import (
	"thyngo/internal/config"

	"github.com/gin-gonic/gin"
)

func NewEngine(cfg config.Config) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	// Middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	//TODO ADD Additional Settings (CORS, security, etc.)

	return r
}
