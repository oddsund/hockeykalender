package ehl

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSeasonUnmarshal(t *testing.T) {
	jsonData := `{
		"uuid": "bir2zwf4qa",
		"names": [
			{"language": "no", "translation": "2025/2026"},
			{"language": "sv", "translation": "2025/2026"},
			{"language": "en", "translation": "2025/2026"}
		]
	}`

	var season Season
	err := json.Unmarshal([]byte(jsonData), &season)
	if err != nil {
		t.Fatalf("failed to unmarshal season: %v", err)
	}

	if season.UUID != "bir2zwf4qa" {
		t.Errorf("expected UUID 'bir2zwf4qa', got '%s'", season.UUID)
	}
	if season.Name != "2025/2026" {
		t.Errorf("expected Name '2025/2026', got '%s'", season.Name)
	}
}

func TestTeamUnmarshal(t *testing.T) {
	jsonData := `{
		"uuid": "qQ0-A0eF1CWG5",
		"code": "VIF",
		"names": {
			"full": "Vålerenga Ishockey Elite",
			"short": "Vålerenga"
		},
		"score": 3,
		"icon": "/images/teams/valerenga.png"
	}`

	var team Team
	err := json.Unmarshal([]byte(jsonData), &team)
	if err != nil {
		t.Fatalf("failed to unmarshal team: %v", err)
	}

	if team.UUID != "qQ0-A0eF1CWG5" {
		t.Errorf("expected UUID 'qQ0-A0eF1CWG5', got '%s'", team.UUID)
	}
	if team.Code != "VIF" {
		t.Errorf("expected Code 'VIF', got '%s'", team.Code)
	}
	if team.FullName != "Vålerenga Ishockey Elite" {
		t.Errorf("expected FullName 'Vålerenga Ishockey Elite', got '%s'", team.FullName)
	}
	if team.ShortName != "Vålerenga" {
		t.Errorf("expected ShortName 'Vålerenga', got '%s'", team.ShortName)
	}
}

func TestGameUnmarshal(t *testing.T) {
	jsonData := `{
		"uuid": "aqdidacsop",
		"rawStartDateTime": "2025-09-11T17:00:00.000Z",
		"state": "pre-game",
		"homeTeam": {
			"uuid": "qQ0-A0eF1CWG5",
			"code": "VIF",
			"names": {
				"full": "Vålerenga Ishockey Elite",
				"short": "Vålerenga"
			},
			"score": 0
		},
		"awayTeam": {
			"uuid": "qQ0-8E8X1CsEP",
			"code": "STH",
			"names": {
				"full": "Storhamar Ishockey Elite",
				"short": "Storhamar"
			},
			"score": 0
		},
		"venueInfo": {
			"uuid": "venue-123",
			"name": "Jordal Amfi"
		}
	}`

	var game Game
	err := json.Unmarshal([]byte(jsonData), &game)
	if err != nil {
		t.Fatalf("failed to unmarshal game: %v", err)
	}

	if game.UUID != "aqdidacsop" {
		t.Errorf("expected UUID 'aqdidacsop', got '%s'", game.UUID)
	}

	expectedTime := time.Date(2025, 9, 11, 17, 0, 0, 0, time.UTC)
	if !game.StartTime.Equal(expectedTime) {
		t.Errorf("expected StartTime '%v', got '%v'", expectedTime, game.StartTime)
	}

	if game.HomeTeam.ShortName != "Vålerenga" {
		t.Errorf("expected HomeTeam.ShortName 'Vålerenga', got '%s'", game.HomeTeam.ShortName)
	}
	if game.AwayTeam.ShortName != "Storhamar" {
		t.Errorf("expected AwayTeam.ShortName 'Storhamar', got '%s'", game.AwayTeam.ShortName)
	}
	if game.Venue != "Jordal Amfi" {
		t.Errorf("expected Venue 'Jordal Amfi', got '%s'", game.Venue)
	}
}

func TestGameSlug(t *testing.T) {
	tests := []struct {
		name     string
		team     Team
		expected string
	}{
		{"Simple name", Team{ShortName: "Vålerenga"}, "valerenga"},
		{"Name with space", Team{ShortName: "Frisk Asker"}, "frisk-asker"},
		{"Name with special chars", Team{ShortName: "Stavanger Oilers"}, "stavanger-oilers"},
		{"Name with Ø", Team{ShortName: "Lørenskog"}, "lorenskog"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.team.Slug()
			if got != tt.expected {
				t.Errorf("Slug() = %s, want %s", got, tt.expected)
			}
		})
	}
}
