package posts

import (
	"context"
	"errors"
	"time"

	"thyngo/internal/database"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgxpool"
)

type pgStore struct {
	ctxTimeout time.Duration
}

// NewPostgresStore returns a PostStore backed by Postgres.
func NewPostgresStore() PostStore {
	return &pgStore{
		ctxTimeout: 5 * time.Second,
	}
}

func (s *pgStore) ListPosts() []Post {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	pool := database.GetPool()
	if pool == nil {
		return nil
	}

	rows, err := pool.Query(ctx, `SELECT slug, title, content::text, created_at, updated_at FROM posts`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var out []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.Slug, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt); err == nil {
			out = append(out, p)
		}
	}
	return out
}

func (s *pgStore) CreatePost(slug, title, content string) (*Post, error) {
	pool := database.GetPool()
	if pool == nil {
		return nil, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	now := time.Now()
	_, err := pool.Exec(ctx, `INSERT INTO posts (slug, title, content, created_at, updated_at) VALUES ($1,$2,$3::jsonb,$4,$5)`, slug, title, content, now, now)
	if err != nil {
		// unique violation handling
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, err
		}
		return nil, err
	}
	return &Post{
		Slug:      slug,
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (s *pgStore) GetPostBySlug(slug string) *Post {
	pool := database.GetPool()
	if pool == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	var p Post
	err := pool.QueryRow(ctx, `SELECT slug, title, content::text, created_at, updated_at FROM posts WHERE slug=$1`, slug).
		Scan(&p.Slug, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		return nil
	}
	return &p
}

func (s *pgStore) UpdatePostBySlug(slug, title, content string) (*Post, error) {
	pool := database.GetPool()
	if pool == nil {
		return nil, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	var p Post
	err := pool.QueryRow(ctx, `UPDATE posts SET title=$1, content=$2::jsonb, updated_at=now() WHERE slug=$3 RETURNING slug, title, content::text, created_at, updated_at`, title, content, slug).
		Scan(&p.Slug, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (s *pgStore) DeletePostBySlug(slug string) (bool, error) {
	pool := database.GetPool()
	if pool == nil {
		return false, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	cmd, err := pool.Exec(ctx, `DELETE FROM posts WHERE slug=$1`, slug)
	if err != nil {
		return false, err
	}
	return cmd.RowsAffected() > 0, nil
}
