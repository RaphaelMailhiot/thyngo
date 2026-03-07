package database

import (
	"context"
	"fmt"
	"log"
	"thyngo/internal/migrations"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// RunMigrations applique toutes les migrations non encore appliquées via goose
func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil {
		return fmt.Errorf("postgres pool is nil")
	}

	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(pool)
	// On ne ferme pas "db" car stdlib le gère comme une surcouche par-dessus le pool existant

	log.Println("Applying database migrations if any...")
	if err := goose.UpContext(ctx, db, "."); err != nil {
		return fmt.Errorf("goose failed to run migrations: %w", err)
	}

	return nil
}
