package resumes

import (
	"context"
	"errors"
	"time"

	"thyngo/internal/database"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Resume struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Job        string    `json:"job"`
	Github     string    `json:"github"`
	Linkedin   string    `json:"linkedin"`
	Website    string    `json:"website"`
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

func (s *Service) ListResumes() []Resume {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	if s.pool == nil {
		return nil
	}

	rows, err := s.pool.Query(ctx, `SELECT id, user_id, name, email, phone, job, github, linkedin, website, visibility, created_at, updated_at FROM resumes`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var out []Resume
	for rows.Next() {
		var r Resume
		if err := rows.Scan(&r.ID, &r.UserID, &r.Name, &r.Email, &r.Phone, &r.Job, &r.Github, &r.Linkedin, &r.Website, &r.Visibility, &r.CreatedAt, &r.UpdatedAt); err == nil {
			out = append(out, r)
		}
	}
	return out
}

func (s *Service) CreateResume(userID int64, name, email, phone, job, github, linkedin, website, visibility string) (*Resume, error) {
	if s.pool == nil {
		return nil, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	if visibility == "" {
		visibility = "public"
	}

	now := time.Now()
	var r Resume
	err := s.pool.QueryRow(ctx, `INSERT INTO resumes (user_id, name, email, phone, job, github, linkedin, website, visibility, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id, user_id, name, email, phone, job, github, linkedin, website, visibility, created_at, updated_at`,
		userID, name, email, phone, job, github, linkedin, website, visibility, now, now).
		Scan(&r.ID, &r.UserID, &r.Name, &r.Email, &r.Phone, &r.Job, &r.Github, &r.Linkedin, &r.Website, &r.Visibility, &r.CreatedAt, &r.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, err
		}
		return nil, err
	}
	return &r, nil
}

func (s *Service) GetResume(id int64) *Resume {
	if s.pool == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	var r Resume
	err := s.pool.QueryRow(ctx, `SELECT id, user_id, name, email, phone, job, github, linkedin, website, visibility, created_at, updated_at FROM resumes WHERE id=$1`, id).
		Scan(&r.ID, &r.UserID, &r.Name, &r.Email, &r.Phone, &r.Job, &r.Github, &r.Linkedin, &r.Website, &r.Visibility, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return nil
	}
	return &r
}

func (s *Service) UpdateResume(userID int64, id int64, name, email, phone, job, github, linkedin, website, visibility string) (*Resume, error) {
	if s.pool == nil {
		return nil, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	if visibility == "" {
		visibility = "public"
	}

	var r Resume
	err := s.pool.QueryRow(ctx, `UPDATE resumes SET name=$1, email=$2, phone=$3, job=$4, github=$5, linkedin=$6, website=$7, visibility=$8, updated_at=now() WHERE id=$9 AND user_id=$10 RETURNING id, user_id, name, email, phone, job, github, linkedin, website, visibility, created_at, updated_at`, name, email, phone, job, github, linkedin, website, visibility, id, userID).
		Scan(&r.ID, &r.UserID, &r.Name, &r.Email, &r.Phone, &r.Job, &r.Github, &r.Linkedin, &r.Website, &r.Visibility, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &r, nil
}

func (s *Service) DeleteResume(userID int64, id int64) (bool, error) {
	if s.pool == nil {
		return false, errors.New("no postgres pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	cmd, err := s.pool.Exec(ctx, `DELETE FROM resumes WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return false, err
	}
	return cmd.RowsAffected() > 0, nil
}
