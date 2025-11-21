package main

import (
	"log"
	"os"
	postsModule "thyngo/internal/modules/posts"

	"thyngo/internal/app"
	//projetsModule "thyngo/internal/modules/projects"
	//resumeModule "thyngo/internal/modules/resume"
)

func main() {
	a := app.NewApp()

	// Register modules
	//a.RegisterModule(projetsModule.New())
	a.RegisterModule(postsModule.New())
	//a.RegisterModule(resumeModule.New())
	//a.RegisterModule(mediaModule.New())

	a.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := a.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
