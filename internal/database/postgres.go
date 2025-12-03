package database

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var PGPool *pgxpool.Pool

// ConnectPG initialise la connexion PostgreSQL depuis POSTGRES_URI (ou valeur par défaut).
func ConnectPG(ctx context.Context) error {
	uri := os.Getenv("POSTGRES_URI")
	if uri == "" {
		uri = "postgres://postgres:postgres@localhost:5432/thyngo_db?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		return err
	}

	// test ping
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return err
	}

	PGPool = pool
	return nil
}

func GetPGPool() *pgxpool.Pool {
	return PGPool
}

func ClosePG(ctx context.Context) {
	if PGPool == nil {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	PGPool.Close()
	PGPool = nil
}
