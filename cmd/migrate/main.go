package main

import (
	"context"
	"log"
	"time"

	"thyngo/internal/database"
	_ "thyngo/internal/migrations"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := database.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = database.Close(context.Background())
	}()

	db := database.GetDatabase("thyngo")
	if db == nil {
		log.Fatal("db is nil")
	}

	if err := database.Run(ctx, db); err != nil {
		log.Fatal("migrations failed:", err)
	}

	log.Println("migrations applied successfully")
}
