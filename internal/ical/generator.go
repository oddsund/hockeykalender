package ical

import (
	"fmt"
	"strings"
	"time"

	"github.com/thomasoddsund/hockeykalender/internal/ehl"
)

const (
	// GameDuration is the default duration of a hockey game
	GameDuration = 2 * time.Hour

	// UID domain for generated events
	UIDDomain = "ehl.hockeykalender"
)

// GenerateCalendar creates an iCal calendar string from games
// teamFilter: if non-empty, only include games involving this team (by ShortName)
// alarms: list of alarms to add to each event
// seasonName: the season name to include in the calendar name
func GenerateCalendar(games []ehl.Game, teamFilter string, alarms []Alarm, seasonName string) string {
	var sb strings.Builder

	// Filter games if team specified
	filteredGames := games
	if teamFilter != "" {
		filteredGames = filterGamesByTeam(games, teamFilter)
	}

	// Calendar header
	sb.WriteString("BEGIN:VCALENDAR\r\n")
	sb.WriteString("VERSION:2.0\r\n")
	sb.WriteString("PRODID:-//Hockeykalender//EHL//NO\r\n")
	sb.WriteString("CALSCALE:GREGORIAN\r\n")
	sb.WriteString("METHOD:PUBLISH\r\n")

	// Calendar name
	calName := "EHL " + seasonName
	if teamFilter != "" {
		calName = teamFilter + " - EHL " + seasonName
	}
	sb.WriteString(fmt.Sprintf("X-WR-CALNAME:%s\r\n", calName))

	// Generate events
	for _, game := range filteredGames {
		sb.WriteString(formatEvent(game, alarms))
	}

	sb.WriteString("END:VCALENDAR\r\n")

	return sb.String()
}

func filterGamesByTeam(games []ehl.Game, teamName string) []ehl.Game {
	var filtered []ehl.Game
	for _, game := range games {
		if game.HomeTeam.ShortName == teamName || game.AwayTeam.ShortName == teamName {
			filtered = append(filtered, game)
		}
	}
	return filtered
}

func formatEvent(game ehl.Game, alarms []Alarm) string {
	var sb strings.Builder

	summary := fmt.Sprintf("%s vs %s", game.HomeTeam.ShortName, game.AwayTeam.ShortName)
	uid := fmt.Sprintf("%s@%s", game.UUID, UIDDomain)
	dtstamp := time.Now().UTC().Format("20060102T150405Z")
	dtstart := game.StartTime.UTC().Format("20060102T150405Z")
	dtend := game.StartTime.Add(GameDuration).UTC().Format("20060102T150405Z")

	sb.WriteString("BEGIN:VEVENT\r\n")
	sb.WriteString(fmt.Sprintf("UID:%s\r\n", uid))
	sb.WriteString(fmt.Sprintf("DTSTAMP:%s\r\n", dtstamp))
	sb.WriteString(fmt.Sprintf("DTSTART:%s\r\n", dtstart))
	sb.WriteString(fmt.Sprintf("DTEND:%s\r\n", dtend))
	sb.WriteString(fmt.Sprintf("SUMMARY:%s\r\n", summary))
	sb.WriteString(fmt.Sprintf("LOCATION:%s\r\n", game.Venue))

	// Add alarms
	for _, alarm := range alarms {
		sb.WriteString(FormatVALARM(alarm, summary))
		sb.WriteString("\r\n")
	}

	sb.WriteString("END:VEVENT\r\n")

	return sb.String()
}
