package posts

import (
	"context"
	"errors"
	"os"
	"time"

	"thyngo/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	Slug      string    `json:"slug" bson:"slug"`
	Title     string    `json:"title" bson:"title"`
	Content   string    `json:"content" bson:"content"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type Service struct {
	coll       *mongo.Collection
	ctxTimeout time.Duration
}

func NewService() PostStore {
	return NewPostgresStore()
}

func (s *Service) ListPosts() []Post {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	if s.coll == nil {
		return nil
	}

	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil
	}
	defer cur.Close(ctx)

	var posts []Post
	for cur.Next(ctx) {
		var p Post
		if err := cur.Decode(&p); err == nil {
			posts = append(posts, p)
		}
	}
	return posts
}

func (s *Service) CreatePost(slug, title, content string) (*Post, error) {
	if s.coll == nil {
		return nil, errors.New("no mongo collection")
	}
	now := time.Now()
	post := &Post{
		Slug:      slug,
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	_, err := s.coll.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *Service) GetPostBySlug(slug string) *Post {
	if s.coll == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	var post Post
	err := s.coll.FindOne(ctx, bson.M{"slug": slug}).Decode(&post)
	if err != nil {
		return nil
	}
	return &post
}

func (s *Service) UpdatePostBySlug(slug, title, content string) (*Post, error) {
	if s.coll == nil {
		return nil, errors.New("no mongo collection")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":      title,
			"content":    content,
			"updated_at": time.Now(),
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated Post
	err := s.coll.FindOneAndUpdate(ctx, bson.M{"slug": slug}, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &updated, nil
}

func (s *Service) DeletePostBySlug(slug string) (bool, error) {
	if s.coll == nil {
		return false, errors.New("no mongo collection")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	res, err := s.coll.DeleteOne(ctx, bson.M{"slug": slug})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}
