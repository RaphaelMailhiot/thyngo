package media

import (
	//"context"
	//"errors"
	"time"

	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
	"thyngo/internal/database"
)

type Media struct {
	ID        string    `json:"id" bson:"id"`
	Slug      string    `json:"slug" bson:"slug"`
	Title     string    `json:"title" bson:"title"`
	Type      string    `json:"type" bson:"type"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type Service struct {
	coll       *mongo.Collection
	ctxTimeout time.Duration
}

func NewService() *Service {
	coll := database.Collection("thyngo", "media")
	return &Service{
		coll:       coll,
		ctxTimeout: 5 * time.Second,
	}
}
