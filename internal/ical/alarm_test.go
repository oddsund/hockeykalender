package ical

import (
	"strings"
	"testing"
	"time"
)

func TestAlarmTrigger(t *testing.T) {
	tests := []struct {
		name     string
		alarm    Alarm
		expected string
	}{
		{"1 day", Alarm1Day, "-P1D"},
		{"3 hours", Alarm3Hours, "-PT3H"},
		{"1 hour", Alarm1Hour, "-PT1H"},
		{"15 minutes", Alarm15Min, "-PT15M"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.alarm.Trigger() != tt.expected {
				t.Errorf("Trigger() = %s, want %s", tt.alarm.Trigger(), tt.expected)
			}
		})
	}
}

func TestAlarmSuffix(t *testing.T) {
	tests := []struct {
		name     string
		alarm    Alarm
		expected string
	}{
		{"1 day", Alarm1Day, "1d"},
		{"3 hours", Alarm3Hours, "3h"},
		{"1 hour", Alarm1Hour, "1h"},
		{"15 minutes", Alarm15Min, "15m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.alarm.Suffix() != tt.expected {
				t.Errorf("Suffix() = %s, want %s", tt.alarm.Suffix(), tt.expected)
			}
		})
	}
}

func TestGenerateAlarmCombinations(t *testing.T) {
	combos := GenerateAlarmCombinations()

	// Should be 16 combinations (2^4)
	if len(combos) != 16 {
		t.Errorf("expected 16 combinations, got %d", len(combos))
	}

	// First should be empty
	if len(combos[0]) != 0 {
		t.Errorf("first combination should be empty, got %d alarms", len(combos[0]))
	}

	// Last should have all 4 alarms
	if len(combos[15]) != 4 {
		t.Errorf("last combination should have 4 alarms, got %d", len(combos[15]))
	}
}

func TestAlarmSetSuffix(t *testing.T) {
	tests := []struct {
		name     string
		alarms   []Alarm
		expected string
	}{
		{"empty", []Alarm{}, ""},
		{"single", []Alarm{Alarm1Day}, "1d"},
		{"two", []Alarm{Alarm1Day, Alarm1Hour}, "1d-1h"},
		{"all", []Alarm{Alarm1Day, Alarm3Hours, Alarm1Hour, Alarm15Min}, "1d-3h-1h-15m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AlarmSetSuffix(tt.alarms)
			if got != tt.expected {
				t.Errorf("AlarmSetSuffix() = %s, want %s", got, tt.expected)
			}
		})
	}
}

func TestFormatVALARM(t *testing.T) {
	alarm := Alarm1Hour
	description := "Vålerenga vs Storhamar"

	result := FormatVALARM(alarm, description)

	if !strings.Contains(result, "BEGIN:VALARM") {
		t.Error("expected BEGIN:VALARM")
	}
	if !strings.Contains(result, "TRIGGER:-PT1H") {
		t.Error("expected TRIGGER:-PT1H")
	}
	if !strings.Contains(result, "ACTION:DISPLAY") {
		t.Error("expected ACTION:DISPLAY")
	}
	if !strings.Contains(result, "DESCRIPTION:Vålerenga vs Storhamar") {
		t.Error("expected DESCRIPTION with match info")
	}
	if !strings.Contains(result, "END:VALARM") {
		t.Error("expected END:VALARM")
	}
}

func TestAlarmDescription(t *testing.T) {
	tests := []struct {
		alarm    Alarm
		match    string
		expected string
	}{
		{Alarm1Day, "Vålerenga vs Storhamar", "Vålerenga vs Storhamar i morgen"},
		{Alarm3Hours, "Vålerenga vs Storhamar", "Vålerenga vs Storhamar om 3 timer"},
		{Alarm1Hour, "Vålerenga vs Storhamar", "Vålerenga vs Storhamar om 1 time"},
		{Alarm15Min, "Vålerenga vs Storhamar", "Vålerenga vs Storhamar om 15 minutter"},
	}

	for _, tt := range tests {
		t.Run(tt.alarm.Suffix(), func(t *testing.T) {
			got := tt.alarm.Description(tt.match)
			if got != tt.expected {
				t.Errorf("Description() = %s, want %s", got, tt.expected)
			}
		})
	}
}

func TestAlarmDuration(t *testing.T) {
	tests := []struct {
		alarm    Alarm
		expected time.Duration
	}{
		{Alarm1Day, 24 * time.Hour},
		{Alarm3Hours, 3 * time.Hour},
		{Alarm1Hour, 1 * time.Hour},
		{Alarm15Min, 15 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.alarm.Suffix(), func(t *testing.T) {
			got := tt.alarm.Duration()
			if got != tt.expected {
				t.Errorf("Duration() = %v, want %v", got, tt.expected)
			}
		})
	}
}
