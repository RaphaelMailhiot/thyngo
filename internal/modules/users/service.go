package users

import (
	"context"
	"errors"
	"time"

	"thyngo/internal/auth"
	"thyngo/internal/database"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Don't expose in JSON
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
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

// RegisterUser creates a new user, hashes the password, and saves to DB.
func (s *Service) RegisterUser(username, email, password string) (*User, error) {
	if s.pool == nil {
		return nil, errors.New("database pool is not initialized")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	now := time.Now()
	defaultRole := "user"
	var u User
	err = s.pool.QueryRow(ctx, `INSERT INTO users (username, email, password_hash, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, username, email, role, created_at, updated_at`,
		username, email, string(hashed), defaultRole, now, now).Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errors.New("username or email already exists")
		}
		return nil, err
	}

	return &u, nil
}

// LoginUser verifies credentials and returns a JWT token.
func (s *Service) LoginUser(identifier, password string) (string, *User, error) {
	if s.pool == nil {
		return "", nil, errors.New("database pool is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	var u User
	// User can login with either email or username
	err := s.pool.QueryRow(ctx, `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1 OR username = $1`, identifier).
		Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil, errors.New("invalid credentials")
		}
		return "", nil, err
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Generate JWT Token
	token, err := auth.GenerateToken(u.ID, u.Role)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	return token, &u, nil
}
