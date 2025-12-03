package database

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

// Connect établit la connexion au pool Postgres.
// Lit la chaîne de connexion depuis l'env `DATABASE_URL`.
// Retourne une erreur si la connexion ou le test échoue.
func Connect(ctx context.Context) error {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// valeur par défaut pratique pour dev local
		connStr = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	// Optionnel : créer un contexte avec timeout pour l'initialisation
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	p, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return err
	}

	// simple vérification de santé
	var one int
	if err := p.QueryRow(ctx, "SELECT 1").Scan(&one); err != nil {
		p.Close()
		return err
	}
	if one != 1 {
		p.Close()
		return errors.New("postgres health check failed")
	}

	pool = p
	return nil
}

// Close ferme le pool si présent.
func Close(_ context.Context) error {
	if pool != nil {
		pool.Close()
		pool = nil
	}
	return nil
}

// GetPool retourne le pool actif (ou nil).
func GetPool() *pgxpool.Pool {
	return pool
}
