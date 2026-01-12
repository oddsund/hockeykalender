package ical

import (
	"fmt"
	"strings"
	"time"
)

// Alarm represents a reminder before an event
type Alarm int

const (
	Alarm1Day Alarm = iota
	Alarm3Hours
	Alarm1Hour
	Alarm15Min
)

// AllAlarms contains all available alarm types in order
var AllAlarms = []Alarm{Alarm1Day, Alarm3Hours, Alarm1Hour, Alarm15Min}

// Trigger returns the iCal TRIGGER value for this alarm
func (a Alarm) Trigger() string {
	switch a {
	case Alarm1Day:
		return "-P1D"
	case Alarm3Hours:
		return "-PT3H"
	case Alarm1Hour:
		return "-PT1H"
	case Alarm15Min:
		return "-PT15M"
	default:
		return ""
	}
}

// Suffix returns the filename suffix for this alarm
func (a Alarm) Suffix() string {
	switch a {
	case Alarm1Day:
		return "1d"
	case Alarm3Hours:
		return "3h"
	case Alarm1Hour:
		return "1h"
	case Alarm15Min:
		return "15m"
	default:
		return ""
	}
}

// Duration returns the time.Duration for this alarm
func (a Alarm) Duration() time.Duration {
	switch a {
	case Alarm1Day:
		return 24 * time.Hour
	case Alarm3Hours:
		return 3 * time.Hour
	case Alarm1Hour:
		return 1 * time.Hour
	case Alarm15Min:
		return 15 * time.Minute
	default:
		return 0
	}
}

// Description returns a Norwegian description for the alarm
func (a Alarm) Description(matchSummary string) string {
	switch a {
	case Alarm1Day:
		return matchSummary + " i morgen"
	case Alarm3Hours:
		return matchSummary + " om 3 timer"
	case Alarm1Hour:
		return matchSummary + " om 1 time"
	case Alarm15Min:
		return matchSummary + " om 15 minutter"
	default:
		return matchSummary
	}
}

// AlarmSetSuffix returns the combined filename suffix for a set of alarms
func AlarmSetSuffix(alarms []Alarm) string {
	if len(alarms) == 0 {
		return ""
	}

	suffixes := make([]string, len(alarms))
	for i, a := range alarms {
		suffixes[i] = a.Suffix()
	}
	return strings.Join(suffixes, "-")
}

// FormatVALARM generates the iCal VALARM component
func FormatVALARM(alarm Alarm, matchSummary string) string {
	return fmt.Sprintf("BEGIN:VALARM\r\nTRIGGER:%s\r\nACTION:DISPLAY\r\nDESCRIPTION:%s\r\nEND:VALARM",
		alarm.Trigger(),
		alarm.Description(matchSummary))
}

// GenerateAlarmCombinations returns all 16 possible combinations of alarms
func GenerateAlarmCombinations() [][]Alarm {
	combinations := make([][]Alarm, 0, 16)

	// Generate all 2^4 = 16 combinations using bit manipulation
	for i := 0; i < 16; i++ {
		var combo []Alarm
		for j, alarm := range AllAlarms {
			if i&(1<<j) != 0 {
				combo = append(combo, alarm)
			}
		}
		combinations = append(combinations, combo)
	}

	return combinations
}
