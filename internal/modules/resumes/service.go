package resumes

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"thyngo/internal/database"
)

type Resume struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Job       string    `json:"job"`
	Github    string    `json:"github"`
	Linkedin  string    `json:"linkedin"`
	Website   string    `json:"website"`
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
