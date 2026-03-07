package media

import (
	"context"
	"errors"
	"time"

	"thyngo/internal/database"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Media struct {
	ID         int64     `json:"id"`
	UserID     *int64    `json:"user_id,omitempty"`
	Slug       string    `json:"slug"`
	Title      string    `json:"title"`
	Type       string    `json:"type"`
	URL        string    `json:"url"`
	FileName   string    `json:"file_name"`
	FileSize   int64     `json:"file_size"`
	MimeType   string    `json:"mime_type"`
	Visibility string    `json:"visibility"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type Service struct {
	pool       *pgxpool.Pool
	ctxTimeout time.Duration
}

func NewService() *Service {
	return &Service{
		pool:       database.GetPool(),
		ctxTimeout: 5 * time.Second,
	}
}

func (s *Service) ListMedia() []Media {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	if s.pool == nil {
		return nil
	}

	rows, err := s.pool.Query(ctx, `SELECT id, user_id, slug, title, type, url, file_name, file_size, mime_type, visibility, created_at, updated_at FROM media`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var out []Media
	for rows.Next() {
		var m Media
		if err := rows.Scan(&m.ID, &m.UserID, &m.Slug, &m.Title, &m.Type, &m.URL, &m.FileName, &m.FileSize, &m.MimeType, &m.Visibility, &m.CreatedAt, &m.UpdatedAt); err == nil {
			out = append(out, m)
		}
	}
	return out
}

func (s *Service) CreateMedia(userID *int64, slug, title, mediaType, url, fileName string, fileSize int64, mimeType, visibility string) (*Media, error) {
	if s.pool == nil {
		return nil, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	if visibility == "" {
		visibility = "public"
	}

	now := time.Now()
	var m Media
	err := s.pool.QueryRow(ctx, `INSERT INTO media (user_id, slug, title, type, url, file_name, file_size, mime_type, visibility, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id, user_id, slug, title, type, url, file_name, file_size, mime_type, visibility, created_at, updated_at`,
		userID, slug, title, mediaType, url, fileName, fileSize, mimeType, visibility, now, now).
		Scan(&m.ID, &m.UserID, &m.Slug, &m.Title, &m.Type, &m.URL, &m.FileName, &m.FileSize, &m.MimeType, &m.Visibility, &m.CreatedAt, &m.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, err
		}
		return nil, err
	}
	return &m, nil
}

func (s *Service) GetMediaBySlug(slug string) *Media {
	if s.pool == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	var m Media
	err := s.pool.QueryRow(ctx, `SELECT id, user_id, slug, title, type, url, file_name, file_size, mime_type, visibility, created_at, updated_at FROM media WHERE slug=$1`, slug).
		Scan(&m.ID, &m.UserID, &m.Slug, &m.Title, &m.Type, &m.URL, &m.FileName, &m.FileSize, &m.MimeType, &m.Visibility, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return nil
	}
	return &m
}

func (s *Service) UpdateMediaBySlug(userID *int64, slug, title, visibility string) (*Media, error) {
	if s.pool == nil {
		return nil, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	if visibility == "" {
		visibility = "public"
	}

	var m Media
	query := `UPDATE media SET title=$1, visibility=$2, updated_at=now() WHERE slug=$3`
	args := []interface{}{title, visibility, slug}

	if userID != nil {
		query += ` AND user_id=$4`
		args = append(args, *userID)
	}

	query += ` RETURNING id, user_id, slug, title, type, url, file_name, file_size, mime_type, visibility, created_at, updated_at`

	err := s.pool.QueryRow(ctx, query, args...).
		Scan(&m.ID, &m.UserID, &m.Slug, &m.Title, &m.Type, &m.URL, &m.FileName, &m.FileSize, &m.MimeType, &m.Visibility, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (s *Service) DeleteMediaBySlug(userID *int64, slug string) (bool, error) {
	if s.pool == nil {
		return false, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	query := `DELETE FROM media WHERE slug=$1`
	args := []interface{}{slug}

	if userID != nil {
		query += ` AND user_id=$2`
		args = append(args, *userID)
	}

	cmd, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		return false, err
	}
	return cmd.RowsAffected() > 0, nil
}
