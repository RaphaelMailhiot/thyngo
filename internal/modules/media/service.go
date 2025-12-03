package media

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"thyngo/internal/database"
)

type Media struct {
	ID        int64     `json:"id"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
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
