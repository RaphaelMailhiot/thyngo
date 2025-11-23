package main

import (
	"context"
	"log"
	"os"
	"thyngo/internal/config"
	"thyngo/internal/database"

	mediaModule "thyngo/internal/modules/media"
	postsModule "thyngo/internal/modules/posts"
	projetsModule "thyngo/internal/modules/projects"
	resumeModule "thyngo/internal/modules/resume"
	usersModule "thyngo/internal/modules/users"

	"thyngo/internal/app"
)

func main() {
	// Connect to MongoDB
	if err := database.Connect(context.Background()); err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	// Close connection when main function ends
	defer func() {
		_ = database.Close(context.Background())
	}()

	// Load feature flags / module toggles
	cfg := config.Load()

	a := app.NewApp()

	// Register modules conditionnellement
	if cfg.EnableMedia {
		a.RegisterModule(mediaModule.New())
	}
	if cfg.EnablePosts {
		a.RegisterModule(postsModule.New())
	}
	if cfg.EnableProjects {
		a.RegisterModule(projetsModule.New())
	}
	if cfg.EnableResume {
		a.RegisterModule(resumeModule.New())
	}
	if cfg.EnableUsers {
		a.RegisterModule(usersModule.New())
	}

	for _, m := range a.Modules {
		log.Printf("Registered module: %s", m.Name())
	}
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
