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
		Module:      "users",
		Description: "Create users collection: unique index on email and username and seed two users",
		Up: func(ctx context.Context, db *mongo.Database) error {
			coll := db.Collection("users")

			// crée des index uniques sur email et username
			indexModels := []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "email", Value: 1}},
					Options: options.Index().SetUnique(true),
				},
				{
					Keys:    bson.D{{Key: "username", Value: 1}},
					Options: options.Index().SetUnique(true),
				},
			}
			if _, err := coll.Indexes().CreateMany(ctx, indexModels); err != nil {
				return err
			}

			// utilisateurs
			users := []struct {
				Username string
				Email    string
				Role     string
			}{
				{Username: "Admin", Email: "admin@example.com", Role: "admin"},
				{Username: "User One", Email: "user1@example.com", Role: "user"},
			}

			for _, u := range users {
				_, err := coll.UpdateOne(
					ctx,
					bson.M{"email": u.Email},
					bson.M{"$setOnInsert": bson.M{
						"username":   u.Username,
						"email":      u.Email,
						"role":       u.Role,
						"created_at": time.Now(),
						"updated_at": time.Now(),
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
