package database

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

// Connect initialise la connexion MongoDB en utilisant MONGO_URI.
func Connect(ctx context.Context) error {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(uri)
	c, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}

	if err := c.Ping(ctx, readpref.Primary()); err != nil {
		_ = c.Disconnect(ctx)
		return err
	}

	client = c
	return nil
}

// GetDatabase retourne une référence à la base demandée.
func GetDatabase(name string) *mongo.Database {
	if client == nil {
		return nil
	}
	return client.Database(name)
}

// Collection retourne une collection pour usage direct.
func Collection(dbName, collName string) *mongo.Collection {
	db := GetDatabase(dbName)
	if db == nil {
		return nil
	}
	return db.Collection(collName)
}

// Close ferme la connexion MongoDB.
func Close(ctx context.Context) error {
	if client == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := client.Disconnect(ctx)
	client = nil
	return err
}
