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

type PostStore interface {
	ListPosts() []Post
	CreatePost(userID *int64, slug, title, visibility string) (*Post, error)
	GetPostBySlug(slug string) *Post
	GetPostBySlugWithContents(slug string) *Post
	UpdatePostBySlug(slug, title, visibility string) (*Post, error)
	DeletePostBySlug(slug string) (bool, error)

	// Block management
	CreateBlock(blockType string, data BlockData) (*Block, error)
	GetBlock(blockID int64) *Block
	UpdateBlock(blockID int64, data BlockData) (*Block, error)
	DeleteBlock(blockID int64) (bool, error)

	// Content management
	CreateContent(parentID int64, order int, blockType string, blockID int64) (*Content, error)
	ListContentsByPostID(postID int64) []Content
	UpdateContent(contentID int64, order int) (*Content, error)
	DeleteContent(contentID int64) (bool, error)
}

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

	rows, err := pool.Query(ctx, `SELECT id, user_id, slug, title, visibility, created_at, updated_at FROM posts`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var out []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Slug, &p.Title, &p.Visibility, &p.CreatedAt, &p.UpdatedAt); err == nil {
			out = append(out, p)
		}
	}
	return out
}

func (s *pgStore) CreatePost(userID *int64, slug, title, visibility string) (*Post, error) {
	pool := database.GetPool()
	if pool == nil {
		return nil, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	if visibility == "" {
		visibility = "public"
	}

	now := time.Now()
	var p Post
	err := pool.QueryRow(ctx, `INSERT INTO posts (user_id, slug, title, visibility, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, user_id, slug, title, visibility, created_at, updated_at`,
		userID, slug, title, visibility, now, now).
		Scan(&p.ID, &p.UserID, &p.Slug, &p.Title, &p.Visibility, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		// unique violation handling
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, err
		}
		return nil, err
	}
	return &p, nil
}

func (s *pgStore) GetPostBySlug(slug string) *Post {
	pool := database.GetPool()
	if pool == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	var p Post
	err := pool.QueryRow(ctx, `SELECT id, user_id, slug, title, visibility, created_at, updated_at FROM posts WHERE slug=$1`, slug).
		Scan(&p.ID, &p.UserID, &p.Slug, &p.Title, &p.Visibility, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		return nil
	}
	return &p
}

func (s *pgStore) GetPostBySlugWithContents(slug string) *Post {
	p := s.GetPostBySlug(slug)
	if p == nil {
		return nil
	}
	p.Contents = s.ListContentsByPostID(p.ID)
	return p
}

func (s *pgStore) UpdatePostBySlug(slug, title, visibility string) (*Post, error) {
	pool := database.GetPool()
	if pool == nil {
		return nil, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	if visibility == "" {
		visibility = "public"
	}

	var p Post
	err := pool.QueryRow(ctx, `UPDATE posts SET title=$1, visibility=$2, updated_at=now() WHERE slug=$3 RETURNING id, user_id, slug, title, visibility, created_at, updated_at`, title, visibility, slug).
		Scan(&p.ID, &p.UserID, &p.Slug, &p.Title, &p.Visibility, &p.CreatedAt, &p.UpdatedAt)
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

// Block management

func (s *pgStore) CreateBlock(blockType string, data BlockData) (*Block, error) {
	pool := database.GetPool()
	if pool == nil {
		return nil, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	now := time.Now()
	var b Block
	err := pool.QueryRow(ctx, `INSERT INTO blocks (type, data, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, type, data, created_at, updated_at`, blockType, data, now, now).
		Scan(&b.ID, &b.Type, &b.Data, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (s *pgStore) GetBlock(blockID int64) *Block {
	pool := database.GetPool()
	if pool == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	var b Block
	err := pool.QueryRow(ctx, `SELECT id, type, data, created_at, updated_at FROM blocks WHERE id=$1`, blockID).
		Scan(&b.ID, &b.Type, &b.Data, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		return nil
	}
	return &b
}

func (s *pgStore) UpdateBlock(blockID int64, data BlockData) (*Block, error) {
	pool := database.GetPool()
	if pool == nil {
		return nil, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	var b Block
	err := pool.QueryRow(ctx, `UPDATE blocks SET data=$1, updated_at=now() WHERE id=$2 RETURNING id, type, data, created_at, updated_at`, data, blockID).
		Scan(&b.ID, &b.Type, &b.Data, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

func (s *pgStore) DeleteBlock(blockID int64) (bool, error) {
	pool := database.GetPool()
	if pool == nil {
		return false, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	cmd, err := pool.Exec(ctx, `DELETE FROM blocks WHERE id=$1`, blockID)
	if err != nil {
		return false, err
	}
	return cmd.RowsAffected() > 0, nil
}

// Content management

func (s *pgStore) CreateContent(parentID int64, order int, blockType string, blockID int64) (*Content, error) {
	pool := database.GetPool()
	if pool == nil {
		return nil, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	now := time.Now()
	var c Content
	err := pool.QueryRow(ctx, `INSERT INTO contents (parent_table, parent_id, ord, type, block_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, parent_id, ord, type, block_id, created_at, updated_at`,
		"posts", parentID, order, blockType, blockID, now, now).
		Scan(&c.ID, &c.ParentID, &c.Order, &c.Type, &blockID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	c.Block = s.GetBlock(blockID)
	return &c, nil
}

func (s *pgStore) ListContentsByPostID(postID int64) []Content {
	pool := database.GetPool()
	if pool == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	rows, err := pool.Query(ctx, `SELECT id, parent_id, ord, type, block_id, created_at, updated_at FROM contents WHERE parent_table=$1 AND parent_id=$2 ORDER BY ord`, "posts", postID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var out []Content
	for rows.Next() {
		var c Content
		var blockID int64
		if err := rows.Scan(&c.ID, &c.ParentID, &c.Order, &c.Type, &blockID, &c.CreatedAt, &c.UpdatedAt); err == nil {
			c.Block = s.GetBlock(blockID)
			out = append(out, c)
		}
	}
	return out
}

func (s *pgStore) UpdateContent(contentID int64, order int) (*Content, error) {
	pool := database.GetPool()
	if pool == nil {
		return nil, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	var c Content
	var blockID int64
	err := pool.QueryRow(ctx, `UPDATE contents SET ord=$1, updated_at=now() WHERE id=$2 RETURNING id, parent_id, ord, type, block_id, created_at, updated_at`, order, contentID).
		Scan(&c.ID, &c.ParentID, &c.Order, &c.Type, &blockID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	c.Block = s.GetBlock(blockID)
	return &c, nil
}

func (s *pgStore) DeleteContent(contentID int64) (bool, error) {
	pool := database.GetPool()
	if pool == nil {
		return false, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	cmd, err := pool.Exec(ctx, `DELETE FROM contents WHERE id=$1`, contentID)
	if err != nil {
		return false, err
	}
	return cmd.RowsAffected() > 0, nil
}
