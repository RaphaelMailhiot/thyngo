package media

import (
	"encoding/json"
	"fmt"
	"net/http"
	"thyngo/internal/config"
)

type FanartClient struct {
	apiKey  string
	baseURL string
}

func NewFanartClient() *FanartClient {
	cfg := config.Load()
	return &FanartClient{
		apiKey:  cfg.FanartAPIKey,
		baseURL: "https://webservice.fanart.tv/v3",
	}
}

type FanartMusicResponse struct {
	Albums map[string]struct {
		AlbumCover []struct {
			ID   string `json:"id"`
			URL  string `json:"url"`
			Likes string `json:"likes"`
		} `json:"albumcover"`
	} `json:"albums"`
}

func (c *FanartClient) GetAlbumArtwork(mbid string) ([]string, error) {
	u := fmt.Sprintf("%s/music/albums/%s?api_key=%s", c.baseURL, mbid, c.apiKey)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fanart api returned status: %d", resp.StatusCode)
	}

	var result FanartMusicResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var urls []string
	for _, album := range result.Albums {
		for _, cover := range album.AlbumCover {
			urls = append(urls, cover.URL)
		}
	}

	return urls, nil
}

type FanartMovieResponse struct {
	MoviePoster []struct {
		ID   string `json:"id"`
		URL  string `json:"url"`
		Likes string `json:"likes"`
	} `json:"movieposter"`
}

func (c *FanartClient) GetMovieArtwork(tmdbID string) ([]string, error) {
	u := fmt.Sprintf("%s/movies/%s?api_key=%s", c.baseURL, tmdbID, c.apiKey)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fanart api returned status: %d", resp.StatusCode)
	}

	var result FanartMovieResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var urls []string
	for _, poster := range result.MoviePoster {
		urls = append(urls, poster.URL)
	}

	return urls, nil
}

type FanartSeriesResponse struct {
	TVPoster []struct {
		ID   string `json:"id"`
		URL  string `json:"url"`
		Likes string `json:"likes"`
	} `json:"tvposter"`
}

func (c *FanartClient) GetSeriesArtwork(tmdbID string) ([]string, error) {
	u := fmt.Sprintf("%s/tv/%s?api_key=%s", c.baseURL, tmdbID, c.apiKey)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fanart api returned status: %d", resp.StatusCode)
	}

	var result FanartSeriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var urls []string
	for _, poster := range result.TVPoster {
		urls = append(urls, poster.URL)
	}

	return urls, nil
}
