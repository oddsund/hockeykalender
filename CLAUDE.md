# Claude Code Instructions

This file provides context for Claude Code when working on this project.

## Project Overview

EHL Kalender generates iCal calendar files for EliteHockey Ligaen (Norwegian ice hockey). It fetches game data from the EHL API, generates 176 calendar files (11 entities × 16 alarm combinations), and deploys to GitHub Pages.

## Tech Stack

- **Language:** Go 1.22+
- **External dependencies:** `golang.org/x/text` (Unicode normalization for slugs)
- **Deployment:** GitHub Pages via GitHub Actions
- **Version control:** jujutsu (jj) - ignore `.jj/` directory

## Architecture

```
cmd/generate/main.go     → CLI entrypoint, orchestrates the pipeline
internal/ehl/            → API client + data models (Season, Game, Team)
internal/ical/           → iCal generation (calendar, events, alarms)
internal/output/         → File writing, filename generation
web/                     → Static landing page (copied to dist/)
```

## Key Design Decisions

1. **No external iCal library** - iCal format is simple enough to generate manually
2. **TDD approach** - All packages have comprehensive tests
3. **Flexible JSON parsing** - Types handle both test data format and actual API format
4. **Norwegian user-facing content** - UI and alarm descriptions are in Norwegian

## EHL API

**Season endpoint:**
```
GET https://www.ehl.no/api/sports-v2/season-series-game-types-filter?series=qUu-397s1Dpwm
Response: { "season": [...], "series": [...], "gameType": [...] }
```

**Games endpoint:**
```
GET https://www.ehl.no/api/sports-v2/game-schedule?seasonUuid={uuid}&seriesUuid=qUu-397s1Dpwm&gameTypeUuid=qQ9-af37Ti40B&gamePlace=all&played=all
Response: { "gameInfo": [...] }
```

**Fixed UUIDs:**
- Series (EHL): `qUu-397s1Dpwm`
- Game type (regular season): `qQ9-af37Ti40B`

## Common Tasks

### Run tests
```bash
go test ./...
```

### Generate calendars locally
```bash
go build -o bin/generate ./cmd/generate
./bin/generate -output dist
```

### Add a new alarm option
1. Add constant in `internal/ical/alarm.go`
2. Update `AllAlarms` slice
3. Update `Trigger()`, `Suffix()`, `Duration()`, `Description()` methods
4. Add checkbox in `web/index.html`
5. Update tests

### Update for new season
The app automatically detects the current season from the API (first in list).

## Code Style

- English for code and comments
- Norwegian for user-facing strings (UI, alarm descriptions)
- Use table-driven tests
- Keep packages focused and small

## Testing Notes

- `internal/ehl/client_test.go` uses `httptest.Server` for mocking
- Test data in client tests uses `homeTeamInfo`/`awayTeamInfo` format (matches API)
- `internal/ical/generator_test.go` creates test games manually (not via JSON)

## Files to Never Edit

- `.jj/` - jujutsu version control
- `bin/` - compiled binaries (gitignored)
- `dist/` - generated output (gitignored)
- `go.sum` - auto-managed by Go
