package output

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/thomasoddsund/hockeykalender/internal/ehl"
	"github.com/thomasoddsund/hockeykalender/internal/ical"
)

// Stats contains statistics about generated files
type Stats struct {
	FilesWritten int
	TotalBytes   int64
}

// Filename generates the filename for a calendar file
func Filename(slug string, alarms []ical.Alarm) string {
	suffix := ical.AlarmSetSuffix(alarms)
	if suffix == "" {
		return slug + ".ics"
	}
	return fmt.Sprintf("%s-%s.ics", slug, suffix)
}

// WriteCalendar writes a calendar string to a file
func WriteCalendar(dir, filename, content string) error {
	path := filepath.Join(dir, filename)
	return os.WriteFile(path, []byte(content), 0644)
}

// GenerateAllCalendars generates all calendar files (teams + EHL, all alarm combos)
func GenerateAllCalendars(dir string, games []ehl.Game, teams []ehl.Team, seasonName string) (Stats, error) {
	var stats Stats

	// Ensure output directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return stats, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Get all alarm combinations
	alarmCombos := ical.GenerateAlarmCombinations()

	// Generate files for each team
	for _, team := range teams {
		slug := team.Slug()

		for _, alarms := range alarmCombos {
			content := ical.GenerateCalendar(games, team.ShortName, alarms, seasonName)
			filename := Filename(slug, alarms)

			if err := WriteCalendar(dir, filename, content); err != nil {
				return stats, fmt.Errorf("failed to write %s: %w", filename, err)
			}

			stats.FilesWritten++
			stats.TotalBytes += int64(len(content))
		}
	}

	// Generate files for all EHL games
	for _, alarms := range alarmCombos {
		content := ical.GenerateCalendar(games, "", alarms, seasonName)
		filename := Filename("ehl", alarms)

		if err := WriteCalendar(dir, filename, content); err != nil {
			return stats, fmt.Errorf("failed to write %s: %w", filename, err)
		}

		stats.FilesWritten++
		stats.TotalBytes += int64(len(content))
	}

	return stats, nil
}
