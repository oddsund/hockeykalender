package ical

import (
	"strings"
	"testing"
	"time"

	"github.com/thomasoddsund/hockeykalender/internal/ehl"
)

func makeTestGames() []ehl.Game {
	return []ehl.Game{
		{
			UUID:      "game-1",
			StartTime: time.Date(2025, 9, 11, 17, 0, 0, 0, time.UTC),
			State:     "pre-game",
			HomeTeam: ehl.Team{
				UUID:      "team-vif",
				Code:      "VIF",
				ShortName: "Vålerenga",
			},
			AwayTeam: ehl.Team{
				UUID:      "team-sth",
				Code:      "STH",
				ShortName: "Storhamar",
			},
			Venue: "Jordal Amfi",
		},
		{
			UUID:      "game-2",
			StartTime: time.Date(2025, 9, 12, 18, 0, 0, 0, time.UTC),
			State:     "pre-game",
			HomeTeam: ehl.Team{
				UUID:      "team-sth",
				Code:      "STH",
				ShortName: "Storhamar",
			},
			AwayTeam: ehl.Team{
				UUID:      "team-vif",
				Code:      "VIF",
				ShortName: "Vålerenga",
			},
			Venue: "CC Amfi",
		},
		{
			UUID:      "game-3",
			StartTime: time.Date(2025, 9, 13, 19, 0, 0, 0, time.UTC),
			State:     "pre-game",
			HomeTeam: ehl.Team{
				UUID:      "team-fri",
				Code:      "FRI",
				ShortName: "Frisk Asker",
			},
			AwayTeam: ehl.Team{
				UUID:      "team-oil",
				Code:      "OIL",
				ShortName: "Stavanger Oilers",
			},
			Venue: "Askerhallen",
		},
	}
}

func TestGenerateCalendar_BasicStructure(t *testing.T) {
	games := makeTestGames()

	result := GenerateCalendar(games, "", []Alarm{}, "2025/2026")

	// Check header
	if !strings.HasPrefix(result, "BEGIN:VCALENDAR") {
		t.Error("expected calendar to start with BEGIN:VCALENDAR")
	}
	if !strings.Contains(result, "VERSION:2.0") {
		t.Error("expected VERSION:2.0")
	}
	if !strings.Contains(result, "PRODID:-//Hockeykalender//EHL//NO") {
		t.Error("expected PRODID")
	}
	if !strings.HasSuffix(strings.TrimSpace(result), "END:VCALENDAR") {
		t.Error("expected calendar to end with END:VCALENDAR")
	}
}

func TestGenerateCalendar_AllGames(t *testing.T) {
	games := makeTestGames()

	result := GenerateCalendar(games, "", []Alarm{}, "2025/2026")

	// Should contain all 3 games
	count := strings.Count(result, "BEGIN:VEVENT")
	if count != 3 {
		t.Errorf("expected 3 events, got %d", count)
	}

	// Check calendar name for all games
	if !strings.Contains(result, "X-WR-CALNAME:EHL 2025/2026") {
		t.Error("expected X-WR-CALNAME:EHL 2025/2026")
	}
}

func TestGenerateCalendar_FilterByTeam(t *testing.T) {
	games := makeTestGames()

	result := GenerateCalendar(games, "Vålerenga", []Alarm{}, "2025/2026")

	// Should only contain games involving Vålerenga (2 games)
	count := strings.Count(result, "BEGIN:VEVENT")
	if count != 2 {
		t.Errorf("expected 2 events for Vålerenga, got %d", count)
	}

	// Check calendar name for team
	if !strings.Contains(result, "X-WR-CALNAME:Vålerenga - EHL 2025/2026") {
		t.Error("expected X-WR-CALNAME:Vålerenga - EHL 2025/2026")
	}
}

func TestGenerateCalendar_EventContent(t *testing.T) {
	games := makeTestGames()[:1] // Just first game

	result := GenerateCalendar(games, "", []Alarm{}, "2025/2026")

	// Check UID format
	if !strings.Contains(result, "UID:game-1@ehl.hockeykalender") {
		t.Error("expected UID:game-1@ehl.hockeykalender")
	}

	// Check DTSTART (UTC format)
	if !strings.Contains(result, "DTSTART:20250911T170000Z") {
		t.Error("expected DTSTART:20250911T170000Z")
	}

	// Check DTEND (2 hours later)
	if !strings.Contains(result, "DTEND:20250911T190000Z") {
		t.Error("expected DTEND:20250911T190000Z")
	}

	// Check SUMMARY
	if !strings.Contains(result, "SUMMARY:Vålerenga vs Storhamar") {
		t.Error("expected SUMMARY:Vålerenga vs Storhamar")
	}

	// Check LOCATION
	if !strings.Contains(result, "LOCATION:Jordal Amfi") {
		t.Error("expected LOCATION:Jordal Amfi")
	}
}

func TestGenerateCalendar_WithAlarms(t *testing.T) {
	games := makeTestGames()[:1]

	result := GenerateCalendar(games, "", []Alarm{Alarm1Day, Alarm1Hour}, "2025/2026")

	// Should have 2 alarms
	alarmCount := strings.Count(result, "BEGIN:VALARM")
	if alarmCount != 2 {
		t.Errorf("expected 2 alarms, got %d", alarmCount)
	}

	if !strings.Contains(result, "TRIGGER:-P1D") {
		t.Error("expected 1 day alarm trigger")
	}
	if !strings.Contains(result, "TRIGGER:-PT1H") {
		t.Error("expected 1 hour alarm trigger")
	}
}

func TestGenerateCalendar_NoAlarms(t *testing.T) {
	games := makeTestGames()[:1]

	result := GenerateCalendar(games, "", []Alarm{}, "2025/2026")

	// Should have no alarms
	if strings.Contains(result, "BEGIN:VALARM") {
		t.Error("expected no alarms")
	}
}

func TestGenerateCalendar_LineEndings(t *testing.T) {
	games := makeTestGames()[:1]

	result := GenerateCalendar(games, "", []Alarm{}, "2025/2026")

	// iCal spec requires CRLF line endings
	if !strings.Contains(result, "\r\n") {
		t.Error("expected CRLF line endings")
	}
}
