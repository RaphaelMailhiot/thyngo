package media

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

type Movie struct {
	ID          int64     `json:"id"`
	UserID      *int64    `json:"user_id,omitempty"`
	MediaID     int64     `json:"media_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
	Duration    int       `json:"duration,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type Series struct {
	ID          int64     `json:"id"`
	UserID      *int64    `json:"user_id,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type Season struct {
	ID           int64     `json:"id"`
	SeriesID     int64     `json:"series_id"`
	SeasonNumber int       `json:"season_number"`
	Title        string    `json:"title,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

type Episode struct {
	ID            int64     `json:"id"`
	SeasonID      int64     `json:"season_id"`
	MediaID       int64     `json:"media_id"`
	EpisodeNumber int       `json:"episode_number"`
	Title         string    `json:"title,omitempty"`
	Description   string    `json:"description,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

type MusicAlbum struct {
	ID           int64     `json:"id"`
	UserID       *int64    `json:"user_id,omitempty"`
	Title        string    `json:"title"`
	Artist       string    `json:"artist,omitempty"`
	CoverMediaID *int64    `json:"cover_media_id,omitempty"`
	ReleaseDate  *time.Time `json:"release_date,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

type MusicTrack struct {
	ID          int64     `json:"id"`
	AlbumID     int64     `json:"album_id"`
	MediaID     int64     `json:"media_id"`
	Title       string    `json:"title"`
	TrackNumber int       `json:"track_number,omitempty"`
	Duration    int       `json:"duration,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
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

func (s *Service) ImportMediaFromURL(userID *int64, title, externalURL, visibility string) (*Media, error) {
	if s.pool == nil {
		return nil, errors.New("no postgres pool")
	}

	// Fetch the file
	resp, err := http.Get(externalURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external server returned status: %d", resp.StatusCode)
	}

	// Prepare file info
	contentType := resp.Header.Get("Content-Type")
	extension := ".jpg" // Default
	if strings.Contains(contentType, "image/png") {
		extension = ".png"
	} else if strings.Contains(contentType, "image/webp") {
		extension = ".webp"
	}

	// Generate a slug/filename
	slug := fmt.Sprintf("imported-%d", time.Now().UnixNano())
	fileName := slug + extension
	filePath := filepath.Join("uploads", fileName)

	// Ensure uploads directory exists
	if err := os.MkdirAll("uploads", 0755); err != nil {
		return nil, fmt.Errorf("failed to create uploads dir: %w", err)
	}

	// Save to local disk
	out, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create local file: %w", err)
	}
	defer out.Close()

	size, err := io.Copy(out, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to save image: %w", err)
	}

	// Create media record
	url := "/" + strings.ReplaceAll(filePath, "\\", "/")
	return s.CreateMedia(userID, slug, title, "image", url, fileName, size, contentType, visibility)
}

// Movies

func (s *Service) ListMovies() []Movie {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()
	if s.pool == nil {
		return nil
	}
	rows, _ := s.pool.Query(ctx, `SELECT id, user_id, media_id, title, description, release_date, duration, created_at, updated_at FROM movies`)
	defer rows.Close()
	var out []Movie
	for rows.Next() {
		var m Movie
		rows.Scan(&m.ID, &m.UserID, &m.MediaID, &m.Title, &m.Description, &m.ReleaseDate, &m.Duration, &m.CreatedAt, &m.UpdatedAt)
		out = append(out, m)
	}
	return out
}

func (s *Service) CreateMovie(m *Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()
	return s.pool.QueryRow(ctx, `INSERT INTO movies (user_id, media_id, title, description, release_date, duration) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at, updated_at`,
		m.UserID, m.MediaID, m.Title, m.Description, m.ReleaseDate, m.Duration).
		Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (s *Service) UpdateMovie(m *Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()
	return s.pool.QueryRow(ctx, `UPDATE movies SET title=$1, description=$2, release_date=$3, duration=$4, media_id=$5, updated_at=now() WHERE id=$6 RETURNING updated_at`,
		m.Title, m.Description, m.ReleaseDate, m.Duration, m.MediaID, m.ID).
		Scan(&m.UpdatedAt)
}

// Series

func (s *Service) ListSeries() []Series {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()
	if s.pool == nil {
		return nil
	}
	rows, _ := s.pool.Query(ctx, `SELECT id, user_id, title, description, created_at, updated_at FROM series`)
	defer rows.Close()
	var out []Series
	for rows.Next() {
		var srt Series
		rows.Scan(&srt.ID, &srt.UserID, &srt.Title, &srt.Description, &srt.CreatedAt, &srt.UpdatedAt)
		out = append(out, srt)
	}
	return out
}

func (s *Service) CreateSeries(sr *Series) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()
	return s.pool.QueryRow(ctx, `INSERT INTO series (user_id, title, description) VALUES ($1,$2,$3) RETURNING id, created_at, updated_at`,
		sr.UserID, sr.Title, sr.Description).
		Scan(&sr.ID, &sr.CreatedAt, &sr.UpdatedAt)
}

func (s *Service) UpdateSeries(sr *Series) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()
	return s.pool.QueryRow(ctx, `UPDATE series SET title=$1, description=$2, updated_at=now() WHERE id=$3 RETURNING updated_at`,
		sr.Title, sr.Description, sr.ID).
		Scan(&sr.UpdatedAt)
}

// Music Albums

func (s *Service) ListAlbums() []MusicAlbum {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()
	if s.pool == nil {
		return nil
	}
	rows, _ := s.pool.Query(ctx, `SELECT id, user_id, title, artist, cover_media_id, release_date, created_at, updated_at FROM music_albums`)
	defer rows.Close()
	var out []MusicAlbum
	for rows.Next() {
		var a MusicAlbum
		rows.Scan(&a.ID, &a.UserID, &a.Title, &a.Artist, &a.CoverMediaID, &a.ReleaseDate, &a.CreatedAt, &a.UpdatedAt)
		out = append(out, a)
	}
	return out
}

func (s *Service) CreateAlbum(a *MusicAlbum) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()
	return s.pool.QueryRow(ctx, `INSERT INTO music_albums (user_id, title, artist, cover_media_id, release_date) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at, updated_at`,
		a.UserID, a.Title, a.Artist, a.CoverMediaID, a.ReleaseDate).
		Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (s *Service) UpdateAlbum(a *MusicAlbum) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()
	return s.pool.QueryRow(ctx, `UPDATE music_albums SET title=$1, artist=$2, cover_media_id=$3, release_date=$4, updated_at=now() WHERE id=$5 RETURNING updated_at`,
		a.Title, a.Artist, a.CoverMediaID, a.ReleaseDate, a.ID).
		Scan(&a.UpdatedAt)
}

// Seasons, Episodes, Tracks

func (s *Service) CreateSeason(sn *Season) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()
	return s.pool.QueryRow(ctx, `INSERT INTO seasons (series_id, season_number, title) VALUES ($1,$2,$3) RETURNING id, created_at, updated_at`,
		sn.SeriesID, sn.SeasonNumber, sn.Title).
		Scan(&sn.ID, &sn.CreatedAt, &sn.UpdatedAt)
}

func (s *Service) CreateEpisode(e *Episode) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()
	return s.pool.QueryRow(ctx, `INSERT INTO episodes (season_id, media_id, episode_number, title, description) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at, updated_at`,
		e.SeasonID, e.MediaID, e.EpisodeNumber, e.Title, e.Description).
		Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
}

func (s *Service) CreateTrack(t *MusicTrack) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()
	return s.pool.QueryRow(ctx, `INSERT INTO music_tracks (album_id, media_id, title, track_number, duration) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at, updated_at`,
		t.AlbumID, t.MediaID, t.Title, t.TrackNumber, t.Duration).
		Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (s *Service) ListEpisodesBySeason(seasonID int64) []Episode {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()
	if s.pool == nil {
		return nil
	}
	rows, _ := s.pool.Query(ctx, `SELECT id, season_id, media_id, episode_number, title, description, created_at, updated_at FROM episodes WHERE season_id=$1`, seasonID)
	defer rows.Close()
	var out []Episode
	for rows.Next() {
		var e Episode
		rows.Scan(&e.ID, &e.SeasonID, &e.MediaID, &e.EpisodeNumber, &e.Title, &e.Description, &e.CreatedAt, &e.UpdatedAt)
		out = append(out, e)
	}
	return out
}

func (s *Service) ListTracksByAlbum(albumID int64) []MusicTrack {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()
	if s.pool == nil {
		return nil
	}
	rows, _ := s.pool.Query(ctx, `SELECT id, album_id, media_id, title, track_number, duration, created_at, updated_at FROM music_tracks WHERE album_id=$1`, albumID)
	defer rows.Close()
	var out []MusicTrack
	for rows.Next() {
		var t MusicTrack
		rows.Scan(&t.ID, &t.AlbumID, &t.MediaID, &t.Title, &t.TrackNumber, &t.Duration, &t.CreatedAt, &t.UpdatedAt)
		out = append(out, t)
	}
	return out
}
