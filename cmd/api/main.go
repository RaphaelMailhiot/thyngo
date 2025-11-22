package main

import (
	"log"
	"os"

	//mediaModule "thyngo/internal/modules/media"
	postsModule "thyngo/internal/modules/posts"
	//projetsModule "thyngo/internal/modules/projects"
	//resumeModule "thyngo/internal/modules/resume"
	//usersModule "thyngo/internal/modules/users"

	"thyngo/internal/app"
)

func main() {
	a := app.NewApp()

	// Register modules
	//a.RegisterModule(mediaModule.New())
	a.RegisterModule(postsModule.New())
	//a.RegisterModule(resumeModule.New())
	//a.RegisterModule(projetsModule.New())
	//a.RegisterModule(usersModule.New())

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
