package main

import (
	"context"
	"log"
	"os"
	"thyngo/internal/config"
	"thyngo/internal/database"
	"time"

	mediaModule "thyngo/internal/modules/media"
	postsModule "thyngo/internal/modules/posts"
	projetsModule "thyngo/internal/modules/projects"
	resumesModule "thyngo/internal/modules/resumes"
	usersModule "thyngo/internal/modules/users"

	"thyngo/internal/app"
)

func main() {
	// Connect to PostgreSQL database with retry logic
	maxRetries := 10
	retryDelay := 2 * time.Second
	var err error

	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = database.Connect(ctx)
		cancel()

		if err == nil {
			log.Println("Successfully connected to database")
			break
		}

		if i < maxRetries-1 {
			log.Printf("Failed to connect to database (attempt %d/%d): %v. Retrying in %v...", i+1, maxRetries, err, retryDelay)
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		log.Fatalf("failed to connect to database after %d attempts: %v", maxRetries, err)
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
	if cfg.EnableResumes {
		a.RegisterModule(resumesModule.New())
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
