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
		ID:          "20251122_create_posts",
		Module:      "posts",
		Description: "Create posts collection: unique index on slug and seed two posts",
		Up: func(ctx context.Context, db *mongo.Database) error {
			coll := db.Collection("posts")

			// crée un index unique sur slug
			indexModel := mongo.IndexModel{
				Keys:    bson.D{{Key: "slug", Value: 1}},
				Options: options.Index().SetUnique(true),
			}
			if _, err := coll.Indexes().CreateOne(ctx, indexModel); err != nil {
				return err
			}

			// posts
			posts := []struct {
				ID      string
				Slug    string
				Title   string
				Content string
			}{
				{ID: "1", Slug: "first-post", Title: "First Post", Content: "This is the content of the first post."},
				{ID: "2", Slug: "second-post", Title: "Second Post", Content: "This is the content of the second post."},
			}

			for _, p := range posts {
				_, err := coll.UpdateOne(
					ctx,
					bson.M{"slug": p.Slug},
					bson.M{"$setOnInsert": bson.M{
						"id":         p.ID,
						"slug":       p.Slug,
						"title":      p.Title,
						"content":    p.Content,
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
