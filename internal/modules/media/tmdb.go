package media

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"thyngo/internal/config"
)

type TMDBClient struct {
	apiKey      string
	accessToken string
	baseURL     string
}

func NewTMDBClient() *TMDBClient {
	cfg := config.Load()
	return &TMDBClient{
		apiKey:      cfg.TMDBAPIKey,
		accessToken: cfg.TMDBAccessToken,
		baseURL:     "https://api.themoviedb.org/3",
	}
}

type TMDBMovieSearchResult struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	ReleaseDate string  `json:"release_date"`
	PosterPath  string  `json:"poster_path"`
	VoteAverage float64 `json:"vote_average"`
}

type TMDBMovieSearchResponse struct {
	Results []TMDBMovieSearchResult `json:"results"`
}

func (c *TMDBClient) SearchMovies(query string) ([]TMDBMovieSearchResult, error) {
	u := fmt.Sprintf("%s/search/movie?api_key=%s&query=%s", c.baseURL, c.apiKey, url.QueryEscape(query))
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TMDBMovieSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Results, nil
}

type TMDBSeriesSearchResult struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Overview     string  `json:"overview"`
	FirstAirDate string  `json:"first_air_date"`
	PosterPath   string  `json:"poster_path"`
	VoteAverage  float64 `json:"vote_average"`
}

type TMDBSeriesSearchResponse struct {
	Results []TMDBSeriesSearchResult `json:"results"`
}

func (c *TMDBClient) SearchSeries(query string) ([]TMDBSeriesSearchResult, error) {
	u := fmt.Sprintf("%s/search/tv?api_key=%s&query=%s", c.baseURL, c.apiKey, url.QueryEscape(query))
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TMDBSeriesSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Results, nil
}

func (c *TMDBClient) GetMovieDetails(movieID int64) (*TMDBMovieSearchResult, error) {
	u := fmt.Sprintf("%s/movie/%d?api_key=%s", c.baseURL, movieID, c.apiKey)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TMDBMovieSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *TMDBClient) GetSeriesDetails(seriesID int64) (*TMDBSeriesSearchResult, error) {
	u := fmt.Sprintf("%s/tv/%d?api_key=%s", c.baseURL, seriesID, c.apiKey)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TMDBSeriesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
