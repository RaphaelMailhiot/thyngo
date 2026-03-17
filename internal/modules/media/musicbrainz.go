package media

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type MusicBrainzClient struct {
	baseURL string
}

func NewMusicBrainzClient() *MusicBrainzClient {
	return &MusicBrainzClient{
		baseURL: "https://musicbrainz.org/ws/2",
	}
}

type MBReleaseSearchResult struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Date    string `json:"date"`
	Artist  string `json:"artist"` // Simplified, usually more complex
}

type MBReleaseResponse struct {
	Releases []struct {
		ID           string `json:"id"`
		Title        string `json:"title"`
		Date         string `json:"date"`
		ReleaseGroup struct {
			FirstReleaseDate string `json:"first-release-date"`
		} `json:"release-group"`
		ArtistCredit []struct {
			Name string `json:"name"`
		} `json:"artist-credit"`
	} `json:"releases"`
}

func (c *MusicBrainzClient) SearchAlbums(query string) ([]MBReleaseSearchResult, error) {
	// MusicBrainz requires an User-Agent
	u := fmt.Sprintf("%s/release?query=%s&fmt=json", c.baseURL, url.QueryEscape(query))
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", "ThynGo/1.0 ( https://github.com/RaphaelMailhiot/thyngo )")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result MBReleaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var out []MBReleaseSearchResult
	for _, r := range result.Releases {
		artist := ""
		if len(r.ArtistCredit) > 0 {
			artist = r.ArtistCredit[0].Name
		}
		date := r.Date
		if date == "" {
			date = r.ReleaseGroup.FirstReleaseDate
		}
		out = append(out, MBReleaseSearchResult{
			ID:     r.ID,
			Title:  r.Title,
			Date:   date,
			Artist: artist,
		})
	}

	return out, nil
}
