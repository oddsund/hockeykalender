# EHL Kalender

iCal calendar feeds for EliteHockey Ligaen (Norwegian elite ice hockey league).

Subscribe to your team's schedule directly in your calendar app with configurable reminders.

## Features

- Calendar feeds for all 10 EHL teams + combined league calendar
- 16 alarm configurations per calendar (combinations of 1 day, 3 hours, 1 hour, 15 minutes)
- Automatic daily updates via GitHub Actions
- Simple web UI for selecting team and reminder preferences
- Standard iCal format (RFC 5545) compatible with all major calendar apps

## Usage

### Web Interface

Visit the hosted site and:

1. Select your team
2. Choose reminder preferences
3. Click "Abonner på kalender" to subscribe

### Direct Links

Calendar URLs follow this pattern:

```http
https://<your-domain>/{team}.ics
https://<your-domain>/{team}-{alarms}.ics
```

**Teams:** `frisk-asker`, `lillehammer`, `lorenskog`, `narvik`, `nidaros`, `oilers`, `sparta`, `stjernen`, `storhamar`, `valerenga`, `ehl` (all teams)

**Alarm suffixes:** `1d`, `3h`, `1h`, `15m` (can be combined, e.g., `1d-1h`)

**Examples:**

- `valerenga.ics` - Vålerenga games, no reminders
- `valerenga-1h.ics` - Vålerenga games, 1 hour reminder
- `ehl-1d-1h-15m.ics` - All games, reminders at 1 day, 1 hour, and 15 minutes

## Development

### Prerequisites

- Go 1.22+

### Build & Run

```bash
# Run tests
go test ./...

# Build
go build -o generate ./cmd/generate

# Generate calendars
./generate -output dist
```

### Project Structure

```bash
├── cmd/generate/          # CLI entrypoint
├── internal/
│   ├── ehl/               # EHL API client and data types
│   ├── ical/              # iCal generation
│   └── output/            # File writing utilities
├── web/                   # Landing page (HTML/CSS)
└── .github/workflows/     # GitHub Actions automation
```

### Testing

All packages have unit tests:

```bash
go test ./...
go test -v ./internal/ehl/...    # Verbose output for specific package
go test -cover ./...              # With coverage
```

## Deployment

The project is designed for GitHub Pages deployment:

1. Push to `main` branch
2. Enable GitHub Pages in repository settings (Source: GitHub Actions)
3. Calendars will be available at `https://<user>.github.io/<repo>/`

The workflow runs:

- On every push to `main`
- Daily at 06:00 UTC
- Manually via workflow dispatch

## Data Source

Match data is fetched from the official EHL API at `ehl.no`.

## License

MIT
