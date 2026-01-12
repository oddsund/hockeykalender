package ehl

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const testSeasonsResponse = `{
	"season": [
		{"uuid": "bir2zwf4qa", "names": [{"language": "no", "translation": "2025/2026"}]},
		{"uuid": "qec-2Ioo12KN8s", "names": [{"language": "no", "translation": "2024/2025"}]},
		{"uuid": "qd0-6kFP17sVG", "names": [{"language": "no", "translation": "2023/2024"}]}
	]
}`

const testGamesResponse = `{
	"gameInfo": [
		{
			"uuid": "aqdidacsop",
			"rawStartDateTime": "2025-09-11T17:00:00.000Z",
			"state": "pre-game",
			"homeTeamInfo": {
				"uuid": "qQ0-A0eF1CWG5",
				"code": "VIF",
				"names": {"full": "Vålerenga Ishockey Elite", "short": "Vålerenga"},
				"score": 0
			},
			"awayTeamInfo": {
				"uuid": "qQ0-8E8X1CsEP",
				"code": "STH",
				"names": {"full": "Storhamar Ishockey Elite", "short": "Storhamar"},
				"score": 0
			},
			"venueInfo": {"uuid": "venue-123", "name": "Jordal Amfi"}
		},
		{
			"uuid": "bqdidacsop",
			"rawStartDateTime": "2025-09-12T18:00:00.000Z",
			"state": "pre-game",
			"homeTeamInfo": {
				"uuid": "qQ0-8E8X1CsEP",
				"code": "STH",
				"names": {"full": "Storhamar Ishockey Elite", "short": "Storhamar"},
				"score": 0
			},
			"awayTeamInfo": {
				"uuid": "qQ0-A0eF1CWG5",
				"code": "VIF",
				"names": {"full": "Vålerenga Ishockey Elite", "short": "Vålerenga"},
				"score": 0
			},
			"venueInfo": {"uuid": "venue-456", "name": "CC Amfi"}
		}
	]
}`

func TestFetchSeasons(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/sports-v2/season-series-game-types-filter" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testSeasonsResponse))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	seasons, err := client.FetchSeasons()
	if err != nil {
		t.Fatalf("FetchSeasons failed: %v", err)
	}

	if len(seasons) != 3 {
		t.Errorf("expected 3 seasons, got %d", len(seasons))
	}

	if seasons[0].UUID != "bir2zwf4qa" {
		t.Errorf("expected first season UUID 'bir2zwf4qa', got '%s'", seasons[0].UUID)
	}
	if seasons[0].Name != "2025/2026" {
		t.Errorf("expected first season name '2025/2026', got '%s'", seasons[0].Name)
	}
}

func TestGetCurrentSeason(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testSeasonsResponse))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	season, err := client.GetCurrentSeason()
	if err != nil {
		t.Fatalf("GetCurrentSeason failed: %v", err)
	}

	if season.UUID != "bir2zwf4qa" {
		t.Errorf("expected current season UUID 'bir2zwf4qa', got '%s'", season.UUID)
	}
}

func TestFetchGames(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/sports-v2/game-schedule" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		seasonUUID := r.URL.Query().Get("seasonUuid")
		if seasonUUID != "bir2zwf4qa" {
			t.Errorf("expected seasonUuid 'bir2zwf4qa', got '%s'", seasonUUID)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testGamesResponse))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	games, err := client.FetchGames("bir2zwf4qa")
	if err != nil {
		t.Fatalf("FetchGames failed: %v", err)
	}

	if len(games) != 2 {
		t.Errorf("expected 2 games, got %d", len(games))
	}

	if games[0].UUID != "aqdidacsop" {
		t.Errorf("expected first game UUID 'aqdidacsop', got '%s'", games[0].UUID)
	}
	if games[0].HomeTeam.ShortName != "Vålerenga" {
		t.Errorf("expected home team 'Vålerenga', got '%s'", games[0].HomeTeam.ShortName)
	}
	if games[0].Venue != "Jordal Amfi" {
		t.Errorf("expected venue 'Jordal Amfi', got '%s'", games[0].Venue)
	}
}

func TestExtractTeams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testGamesResponse))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	games, _ := client.FetchGames("bir2zwf4qa")

	teams := ExtractTeams(games)

	if len(teams) != 2 {
		t.Errorf("expected 2 unique teams, got %d", len(teams))
	}

	// Check that both teams are present
	teamNames := make(map[string]bool)
	for _, team := range teams {
		teamNames[team.ShortName] = true
	}

	if !teamNames["Vålerenga"] {
		t.Error("expected Vålerenga in teams")
	}
	if !teamNames["Storhamar"] {
		t.Error("expected Storhamar in teams")
	}
}
