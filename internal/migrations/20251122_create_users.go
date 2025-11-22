package migrations

import (
	"context"
	"time"

	"thyngo/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	database.Register(database.Migration{
		ID:          "20251122_create_users",
		Description: "Create users collection: unique index on email and seed two users",
		Up: func(ctx context.Context, db *mongo.Database) error {
			coll := db.Collection("users")

			// crée un index unique sur email
			indexModel := mongo.IndexModel{
				Keys:    bson.D{{Key: "email", Value: 1}},
				Options: options.Index().SetUnique(true),
			}
			if _, err := coll.Indexes().CreateOne(ctx, indexModel); err != nil {
				return err
			}

			// utilisateurs
			users := []struct {
				ID    string
				Email string
				Name  string
				Role  string
			}{
				{ID: "1", Email: "admin@example.com", Name: "Admin", Role: "admin"},
				{ID: "2", Email: "user1@example.com", Name: "User One", Role: "user"},
			}

			for _, u := range users {
				_, err := coll.UpdateOne(
					ctx,
					bson.M{"email": u.Email},
					bson.M{"$setOnInsert": bson.M{
						"id":         u.ID,
						"email":      u.Email,
						"name":       u.Name,
						"role":       u.Role,
						"created_at": time.Now(),
					}},
					options.Update().SetUpsert(true),
				)
				if err != nil {
					return err
				}
			}

			return nil
		},
	})
}
