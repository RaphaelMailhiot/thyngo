package main

import (
	"context"
	"log"
	"os"
	"strings"
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

	// Lire la liste des modules à migrer depuis MIGRATE_MODULES (virgule\-séparés).
	// Si vide, on prend par défaut "posts".
	modulesEnv := os.Getenv("MIGRATE_MODULES")
	var enabled []string
	if modulesEnv != "" {
		for _, s := range strings.Split(modulesEnv, ",") {
			if m := strings.TrimSpace(s); m != "" {
				enabled = append(enabled, m)
			}
		}
	} else {
		enabled = []string{"posts", "users"}
	}

	if err := database.Run(ctx, db, enabled); err != nil {
		log.Fatal("migrations failed:", err)
	}

	log.Println("migrations applied successfully")
}
