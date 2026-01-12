package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/thomasoddsund/hockeykalender/internal/ehl"
	"github.com/thomasoddsund/hockeykalender/internal/output"
)

func main() {
	outputDir := flag.String("output", "dist", "Output directory for generated files")
	flag.Parse()

	log.Println("Starting EHL calendar generation...")

	// Create API client
	client := ehl.NewDefaultClient()

	// Get current season
	log.Println("Fetching current season...")
	season, err := client.GetCurrentSeason()
	if err != nil {
		log.Fatalf("Failed to get current season: %v", err)
	}
	log.Printf("Current season: %s (UUID: %s)", season.Name, season.UUID)

	// Fetch games
	log.Println("Fetching games...")
	games, err := client.FetchGames(season.UUID)
	if err != nil {
		log.Fatalf("Failed to fetch games: %v", err)
	}
	log.Printf("Found %d games", len(games))

	// Extract teams
	teams := ehl.ExtractTeams(games)
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].ShortName < teams[j].ShortName
	})
	log.Printf("Found %d teams", len(teams))

	for _, team := range teams {
		log.Printf("  - %s (%s)", team.ShortName, team.Slug())
	}

	// Generate calendars
	log.Printf("Generating calendars to %s...", *outputDir)
	stats, err := output.GenerateAllCalendars(*outputDir, games, teams, season.Name)
	if err != nil {
		log.Fatalf("Failed to generate calendars: %v", err)
	}

	log.Printf("Generated %d files (%.2f KB total)", stats.FilesWritten, float64(stats.TotalBytes)/1024)

	// Copy web files to output
	log.Println("Copying web files...")
	if err := copyWebFiles(*outputDir); err != nil {
		log.Fatalf("Failed to copy web files: %v", err)
	}

	// Generate team list for HTML page
	generateTeamList(teams)

	log.Println("Done!")
}

func copyWebFiles(outputDir string) error {
	webFiles := []string{"index.html", "style.css"}

	for _, filename := range webFiles {
		src := filepath.Join("web", filename)
		dst := filepath.Join(outputDir, filename)

		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("failed to copy %s: %w", filename, err)
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func generateTeamList(teams []ehl.Team) {
	fmt.Println("\nTeams for HTML page:")
	for _, team := range teams {
		fmt.Printf(`  {name: "%s", slug: "%s"},`+"\n", team.ShortName, team.Slug())
	}
}

func init() {
	// Ensure we exit with proper code on error
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stderr)
}
