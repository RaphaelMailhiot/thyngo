package projects

import (
	"context"
	"errors"
	"time"

	"thyngo/internal/database"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Project struct {
	ID         int64     `json:"id"`
	UserID     *int64    `json:"user_id,omitempty"`
	Slug       string    `json:"slug"`
	Title      string    `json:"title"`
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

func (s *Service) ListProjects() []Project {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	if s.pool == nil {
		return nil
	}

	rows, err := s.pool.Query(ctx, `SELECT id, user_id, slug, title, visibility, created_at, updated_at FROM projects`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var out []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.UserID, &p.Slug, &p.Title, &p.Visibility, &p.CreatedAt, &p.UpdatedAt); err == nil {
			out = append(out, p)
		}
	}
	return out
}

func (s *Service) CreateProject(userID *int64, slug, title, visibility string) (*Project, error) {
	if s.pool == nil {
		return nil, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	if visibility == "" {
		visibility = "public"
	}

	now := time.Now()
	var p Project
	err := s.pool.QueryRow(ctx, `INSERT INTO projects (user_id, slug, title, visibility, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, user_id, slug, title, visibility, created_at, updated_at`,
		userID, slug, title, visibility, now, now).
		Scan(&p.ID, &p.UserID, &p.Slug, &p.Title, &p.Visibility, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, err
		}
		return nil, err
	}
	return &p, nil
}

func (s *Service) GetProjectBySlug(slug string) *Project {
	if s.pool == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	var p Project
	err := s.pool.QueryRow(ctx, `SELECT id, user_id, slug, title, visibility, created_at, updated_at FROM projects WHERE slug=$1`, slug).
		Scan(&p.ID, &p.UserID, &p.Slug, &p.Title, &p.Visibility, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return nil
	}
	return &p
}

func (s *Service) UpdateProjectBySlug(userID *int64, slug, title, visibility string) (*Project, error) {
	if s.pool == nil {
		return nil, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	if visibility == "" {
		visibility = "public"
	}

	var p Project
	query := `UPDATE projects SET title=$1, visibility=$2, updated_at=now() WHERE slug=$3`
	args := []interface{}{title, visibility, slug}

	if userID != nil {
		query += ` AND user_id=$4`
		args = append(args, *userID)
	}

	query += ` RETURNING id, user_id, slug, title, visibility, created_at, updated_at`

	err := s.pool.QueryRow(ctx, query, args...).
		Scan(&p.ID, &p.UserID, &p.Slug, &p.Title, &p.Visibility, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (s *Service) DeleteProjectBySlug(userID *int64, slug string) (bool, error) {
	if s.pool == nil {
		return false, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	query := `DELETE FROM projects WHERE slug=$1`
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
