package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"thyngo/internal/database"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := database.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = database.Close(context.Background())
	}()

	pool := database.GetPool()
	if pool == nil {
		log.Fatal("pgx pool is nil")
	}

	// Lire modules depuis MIGRATE_MODULES (virgule-séparés). Par défaut posts,users
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

	dir := "internal/migrations"
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal("read migrations dir:", err)
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".sql") {
			files = append(files, name)
		}
	}
	sort.Strings(files)

	for _, fname := range files {
		// filtrer par module (si le nom du fichier contient le nom du module)
		matched := false
		for _, m := range enabled {
			if strings.Contains(fname, m) {
				matched = true
				break
			}
		}
		if !matched {
			continue
		}

		path := filepath.Join(dir, fname)
		b, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("read migration %s: %v", fname, err)
		}
		sql := string(b)
		log.Printf("applying migration %s", fname)

		ctxExec, cancelExec := context.WithTimeout(context.Background(), 15*time.Second)
		_, err = pool.Exec(ctxExec, sql)
		cancelExec()
		if err != nil {
			log.Fatalf("migration %s failed: %v", fname, err)
		}
		log.Printf("applied %s", fname)
	}

	log.Println("migrations applied successfully")
}
