package app

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type Module interface {
	Name() string
	RegisterRoutes(router *gin.RouterGroup)
}

type App struct {
	Engine  *gin.Engine
	Modules []Module
}

func NewApp() *App {
	r := gin.Default()

	return &App{
		Engine:  r,
		Modules: make([]Module, 0),
	}
}

func (a *App) RegisterModule(m Module) {
	a.Modules = append(a.Modules, m)
}

func (a *App) SetupRoutes() {
	// Create uploads directory if it doesn't exist
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		_ = os.Mkdir("uploads", 0755)
	}

	// Serve uploaded files statically
	a.Engine.Static("/uploads", "./uploads")

	api := a.Engine.Group("/api")

	for _, m := range a.Modules {
		log.Printf("Registering module: %s", m.Name())
		group := api.Group("/" + m.Name())
		m.RegisterRoutes(group)
	}
}

func (a *App) Run(addr string) error {
	return a.Engine.Run(addr)
}
