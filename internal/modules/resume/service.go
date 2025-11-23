package resume

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

type Resume struct {
	ID         string    `json:"id" bson:"id"`
	Name       string    `json:"slug" bson:"slug"`
	Email      string    `json:"email" bson:"email"`
	Phone      string    `json:"phone" bson:"phone"`
	Job        string    `json:"job" bson:"job"`
	Github     string    `json:"github" bson:"github"`
	Linkedin   string    `json:"linkedin" bson:"linkedin"`
	Website    string    `json:"website" bson:"website"`
	Experience string    `json:"experience" bson:"experience"`
	Education  string    `json:"education" bson:"education"`
	Skills     string    `json:"skills" bson:"skills"`
	CreatedAt  time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type Service struct {
	coll       *mongo.Collection
	ctxTimeout time.Duration
}

func NewService() *Service {
	coll := database.Collection("thyngo", "resume")
	return &Service{
		coll:       coll,
		ctxTimeout: 5 * time.Second,
	}
}
