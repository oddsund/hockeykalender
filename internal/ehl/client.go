package ehl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// EHL API constants
	DefaultBaseURL = "https://www.ehl.no"
	SeriesUUID     = "qUu-397s1Dpwm"   // EliteHockey Ligaen
	GameTypeUUID   = "qQ9-af37Ti40B"   // Regular season
)

// Client is an HTTP client for the EHL API
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new EHL API client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// NewDefaultClient creates a client with the default EHL base URL
func NewDefaultClient() *Client {
	return NewClient(DefaultBaseURL)
}

// seasonsResponse represents the API response for seasons
type seasonsResponse struct {
	Season []Season `json:"season"`
}

// FetchSeasons retrieves all available seasons from the API
func (c *Client) FetchSeasons() ([]Season, error) {
	endpoint := fmt.Sprintf("%s/api/sports-v2/season-series-game-types-filter?series=%s", c.baseURL, SeriesUUID)

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch seasons: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result seasonsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode seasons response: %w", err)
	}

	return result.Season, nil
}

// GetCurrentSeason returns the most recent season (first in the list)
func (c *Client) GetCurrentSeason() (Season, error) {
	seasons, err := c.FetchSeasons()
	if err != nil {
		return Season{}, err
	}

	if len(seasons) == 0 {
		return Season{}, fmt.Errorf("no seasons found")
	}

	return seasons[0], nil
}

// gamesResponse wraps the API response for games
type gamesResponse struct {
	GameInfo []Game `json:"gameInfo"`
}

// FetchGames retrieves all games for a given season
func (c *Client) FetchGames(seasonUUID string) ([]Game, error) {
	params := url.Values{}
	params.Set("seasonUuid", seasonUUID)
	params.Set("seriesUuid", SeriesUUID)
	params.Set("gameTypeUuid", GameTypeUUID)
	params.Set("gamePlace", "all")
	params.Set("played", "all")

	endpoint := fmt.Sprintf("%s/api/sports-v2/game-schedule?%s", c.baseURL, params.Encode())

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch games: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result gamesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode games response: %w", err)
	}

	return result.GameInfo, nil
}

// ExtractTeams returns a list of unique teams from a list of games
func ExtractTeams(games []Game) []Team {
	seen := make(map[string]Team)

	for _, game := range games {
		if _, exists := seen[game.HomeTeam.UUID]; !exists {
			seen[game.HomeTeam.UUID] = game.HomeTeam
		}
		if _, exists := seen[game.AwayTeam.UUID]; !exists {
			seen[game.AwayTeam.UUID] = game.AwayTeam
		}
	}

	teams := make([]Team, 0, len(seen))
	for _, team := range seen {
		teams = append(teams, team)
	}

	return teams
}
