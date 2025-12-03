package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Migration struct {
	ID          string
	Module      string
	Description string
	Up          func(ctx context.Context, pool *pgxpool.Pool) error
}

var registry []Migration

func Register(m Migration) {
	registry = append(registry, m)
}

// Run applique toutes les migrations non encore appliquées.
// La table `migrations` contient les enregistrements avec `id` = migration ID.
// Si enabledModules est vide, toutes les migrations sont considérées autorisées.
func Run(ctx context.Context, pool *pgxpool.Pool, enabledModules []string) error {
	if pool == nil {
		return fmt.Errorf("postgres pool is nil")
	}

	// Créer la table de suivi des migrations si elle n'existe pas
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS migrations (
			id TEXT PRIMARY KEY,
			description TEXT,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`
	if _, err := pool.Exec(ctx, createTableSQL); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Construire un set des modules autorisés
	allowAll := len(enabledModules) == 0
	allowed := map[string]bool{}
	for _, m := range enabledModules {
		allowed[m] = true
	}

	// Récupère les migrations déjà appliquées
	applied := map[string]bool{}
	rows, err := pool.Query(ctx, "SELECT id FROM migrations ORDER BY applied_at")
	if err != nil {
		return fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan migration id: %w", err)
		}
		applied[id] = true
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Applique les migrations dans l'ordre d'enregistrement
	for _, m := range registry {
		// Filtrer par module si nécessaire
		if !allowAll && m.Module != "" && !allowed[m.Module] {
			continue
		}
		if applied[m.ID] {
			continue
		}
		if err := m.Up(ctx, pool); err != nil {
			return fmt.Errorf("migration %s failed: %w", m.ID, err)
		}
		insertSQL := `INSERT INTO migrations (id, description) VALUES ($1, $2)`
		if _, err := pool.Exec(ctx, insertSQL, m.ID, m.Description); err != nil {
			return fmt.Errorf("recording migration %s failed: %w", m.ID, err)
		}
	}
	return nil
}
