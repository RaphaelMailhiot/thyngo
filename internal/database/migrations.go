package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Migration struct {
	ID          string
	Description string
	Up          func(ctx context.Context, db *mongo.Database) error
}

var registry []Migration

func Register(m Migration) {
	registry = append(registry, m)
}

// Run applique toutes les migrations non encore appliquées.
// La collection `migrations` contient des documents avec `_id` = migration ID.
func Run(ctx context.Context, db *mongo.Database) error {
	if db == nil {
		return fmt.Errorf("database is nil")
	}

	col := db.Collection("migrations")

	// Récupère les migrations déjà appliquées
	applied := map[string]bool{}
	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		// si la collection n'existe pas encore, Find renverra nil/empty sans erreur dans la plupart des cas
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc struct {
			ID string `bson:"_id"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return err
		}
		applied[doc.ID] = true
	}

	// Applique les migrations dans l'ordre d'enregistrement
	for _, m := range registry {
		if applied[m.ID] {
			continue
		}
		if err := m.Up(ctx, db); err != nil {
			return fmt.Errorf("migration %s failed: %w", m.ID, err)
		}
		_, err := col.InsertOne(ctx, bson.M{"_id": m.ID, "description": m.Description, "applied_at": time.Now()})
		if err != nil {
			return fmt.Errorf("recording migration %s failed: %w", m.ID, err)
		}
	}
	return nil
}
