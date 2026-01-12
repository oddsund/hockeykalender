package output

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/thomasoddsund/hockeykalender/internal/ehl"
	"github.com/thomasoddsund/hockeykalender/internal/ical"
)

func TestFilename(t *testing.T) {
	tests := []struct {
		name     string
		slug     string
		alarms   []ical.Alarm
		expected string
	}{
		{"no alarms", "valerenga", []ical.Alarm{}, "valerenga.ics"},
		{"single alarm", "valerenga", []ical.Alarm{ical.Alarm1Day}, "valerenga-1d.ics"},
		{"multiple alarms", "valerenga", []ical.Alarm{ical.Alarm1Day, ical.Alarm1Hour}, "valerenga-1d-1h.ics"},
		{"all alarms", "ehl", []ical.Alarm{ical.Alarm1Day, ical.Alarm3Hours, ical.Alarm1Hour, ical.Alarm15Min}, "ehl-1d-3h-1h-15m.ics"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Filename(tt.slug, tt.alarms)
			if got != tt.expected {
				t.Errorf("Filename() = %s, want %s", got, tt.expected)
			}
		})
	}
}

func TestWriteCalendar(t *testing.T) {
	tmpDir := t.TempDir()

	content := "BEGIN:VCALENDAR\r\nEND:VCALENDAR\r\n"
	filename := "test.ics"

	err := WriteCalendar(tmpDir, filename, content)
	if err != nil {
		t.Fatalf("WriteCalendar failed: %v", err)
	}

	// Check file exists
	path := filepath.Join(tmpDir, filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist at %s", path)
	}

	// Check content
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(data) != content {
		t.Errorf("file content mismatch")
	}
}

func TestGenerateAllCalendars(t *testing.T) {
	tmpDir := t.TempDir()

	teams := []ehl.Team{
		{ShortName: "Vålerenga"},
		{ShortName: "Storhamar"},
	}

	games := []ehl.Game{
		{
			UUID:     "game-1",
			HomeTeam: teams[0],
			AwayTeam: teams[1],
			Venue:    "Test Arena",
		},
	}

	stats, err := GenerateAllCalendars(tmpDir, games, teams, "2025/2026")
	if err != nil {
		t.Fatalf("GenerateAllCalendars failed: %v", err)
	}

	// 3 entities (2 teams + EHL) x 16 alarm combos = 48 files
	expectedFiles := 3 * 16
	if stats.FilesWritten != expectedFiles {
		t.Errorf("expected %d files, got %d", expectedFiles, stats.FilesWritten)
	}

	// Check some files exist
	checkFiles := []string{
		"valerenga.ics",
		"valerenga-1d.ics",
		"storhamar.ics",
		"ehl.ics",
		"ehl-1d-3h-1h-15m.ics",
	}

	for _, f := range checkFiles {
		path := filepath.Join(tmpDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", f)
		}
	}
}
