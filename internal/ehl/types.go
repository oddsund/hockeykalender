package ehl

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Season represents an EHL season
type Season struct {
	UUID string `json:"uuid"`
	Name string `json:"-"`
}

// Translation represents a localized name
type Translation struct {
	Language    string `json:"language"`
	Translation string `json:"translation"`
}

// seasonJSON is used for unmarshaling the nested JSON structure
type seasonJSON struct {
	UUID  string        `json:"uuid"`
	Names []Translation `json:"names"`
}

func (s *Season) UnmarshalJSON(data []byte) error {
	var sj seasonJSON
	if err := json.Unmarshal(data, &sj); err != nil {
		return err
	}
	s.UUID = sj.UUID
	// Get Norwegian name, fallback to first available
	for _, n := range sj.Names {
		if n.Language == "no" {
			s.Name = n.Translation
			break
		}
		if s.Name == "" {
			s.Name = n.Translation
		}
	}
	return nil
}

// Team represents a team in a game
type Team struct {
	UUID      string `json:"uuid"`
	Code      string `json:"code"`
	FullName  string `json:"-"`
	ShortName string `json:"-"`
	Score     int    `json:"score"`
	Icon      string `json:"icon"`
}

// teamJSON is used for unmarshaling the nested JSON structure
type teamJSON struct {
	UUID  string `json:"uuid"`
	Code  string `json:"code"`
	Names struct {
		Full  string `json:"full"`
		Short string `json:"short"`
	} `json:"names"`
	Score json.RawMessage `json:"score"`
	Icon  string          `json:"icon"`
}

func (t *Team) UnmarshalJSON(data []byte) error {
	var tj teamJSON
	if err := json.Unmarshal(data, &tj); err != nil {
		return err
	}
	t.UUID = tj.UUID
	t.Code = tj.Code
	t.FullName = tj.Names.Full
	t.ShortName = tj.Names.Short
	t.Icon = tj.Icon

	// Handle score as either int or string
	if len(tj.Score) > 0 {
		// Try int first
		var scoreInt int
		if err := json.Unmarshal(tj.Score, &scoreInt); err == nil {
			t.Score = scoreInt
		} else {
			// Try string
			var scoreStr string
			if err := json.Unmarshal(tj.Score, &scoreStr); err == nil {
				// Parse string to int, default to 0 on error
				fmt.Sscanf(scoreStr, "%d", &t.Score)
			}
		}
	}
	return nil
}

// teamMarshalJSON is used for marshaling Team to JSON
type teamMarshalJSON struct {
	UUID  string `json:"uuid"`
	Code  string `json:"code"`
	Names struct {
		Full  string `json:"full"`
		Short string `json:"short"`
	} `json:"names"`
	Score int    `json:"score"`
	Icon  string `json:"icon"`
}

// MarshalJSON implements json.Marshaler for Team (used in tests)
func (t Team) MarshalJSON() ([]byte, error) {
	return json.Marshal(teamMarshalJSON{
		UUID:  t.UUID,
		Code:  t.Code,
		Score: t.Score,
		Icon:  t.Icon,
		Names: struct {
			Full  string `json:"full"`
			Short string `json:"short"`
		}{
			Full:  t.FullName,
			Short: t.ShortName,
		},
	})
}

// Slug returns a URL-friendly version of the team's short name
func (t *Team) Slug() string {
	name := strings.ToLower(t.ShortName)

	// Replace Norwegian special characters that don't decompose with NFD
	replacer := strings.NewReplacer(
		"æ", "ae",
		"ø", "o",
	)
	name = replacer.Replace(name)

	// Remove diacritics (å → a, etc.)
	transformer := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(transformer, name)

	// Replace spaces with hyphens
	result = strings.ReplaceAll(result, " ", "-")

	return result
}

// Game represents a single game/match
type Game struct {
	UUID      string    `json:"uuid"`
	StartTime time.Time `json:"-"`
	State     string    `json:"state"`
	HomeTeam  Team      `json:"-"`
	AwayTeam  Team      `json:"-"`
	Venue     string    `json:"-"`
}

// gameJSON is used for unmarshaling the nested JSON structure
// Supports both test format (homeTeam/awayTeam) and API format (homeTeamInfo/awayTeamInfo)
type gameJSON struct {
	UUID             string `json:"uuid"`
	RawStartDateTime string `json:"rawStartDateTime"`
	State            string `json:"state"`
	// Test format
	HomeTeam *Team `json:"homeTeam,omitempty"`
	AwayTeam *Team `json:"awayTeam,omitempty"`
	// API format
	HomeTeamInfo *Team `json:"homeTeamInfo,omitempty"`
	AwayTeamInfo *Team `json:"awayTeamInfo,omitempty"`
	VenueInfo    struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
	} `json:"venueInfo"`
}

func (g *Game) UnmarshalJSON(data []byte) error {
	var gj gameJSON
	if err := json.Unmarshal(data, &gj); err != nil {
		return err
	}

	g.UUID = gj.UUID
	g.State = gj.State
	g.Venue = gj.VenueInfo.Name

	// Support both test format and API format
	if gj.HomeTeamInfo != nil {
		g.HomeTeam = *gj.HomeTeamInfo
	} else if gj.HomeTeam != nil {
		g.HomeTeam = *gj.HomeTeam
	}

	if gj.AwayTeamInfo != nil {
		g.AwayTeam = *gj.AwayTeamInfo
	} else if gj.AwayTeam != nil {
		g.AwayTeam = *gj.AwayTeam
	}

	// Parse the ISO 8601 timestamp
	t, err := time.Parse("2006-01-02T15:04:05.000Z", gj.RawStartDateTime)
	if err != nil {
		return err
	}
	g.StartTime = t

	return nil
}

// InvolvesTeam returns true if the given team (by short name) is playing in this game
func (g *Game) InvolvesTeam(teamShortName string) bool {
	return g.HomeTeam.ShortName == teamShortName || g.AwayTeam.ShortName == teamShortName
}
