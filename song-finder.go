// Program that reads image files from a specified file path,
// then uses Google's ML cloud API for reading text, and
// finally queries Spotify to find matches and create a
// playlist
package main

import (
	"fmt"
	"os"

	"github.com/brozeph/song-finder/internal/repositories"
	"github.com/brozeph/song-finder/internal/services"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ttacon/chalk"
)

type cmdlineOptions struct {
	ImageFilePath string `short:"p" long:"path" description:"Path to image files" required:"true"`
	PlaylistName  string `short:"n" long:"playlist" description:"Name of Spotify playlist to create" required:"true"`
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	var (
		options cmdlineOptions
		parser  = flags.NewParser(&options, flags.Default)
	)

	// parse command line arguments
	if _, err := parser.Parse(); err != nil {
		log.Error().Err(err).Msg("")
		os.Exit(1)
	}

	// scaffold up the app
	screenshotRepository := repositories.NewScreenshotRepository()
	spotifyRepository := repositories.NewSpotifyRepository()
	screenshotService := services.NewScreenshotService(&screenshotRepository, &spotifyRepository)

	// find all of the image files
	screenshots, err := screenshotService.Begin(options.ImageFilePath)

	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Printf(
		"Process completed for %s%d%s files",
		chalk.Blue,
		len(screenshots),
		chalk.Reset)

	for _, ss := range screenshots {
		fmt.Println(chalk.Blue, "File:", chalk.Reset, ss.Path)
		fmt.Println(chalk.Red, "Song:", chalk.Reset, ss.SongSearchTerm)
		fmt.Println(chalk.Green, "Spotify URI:", chalk.Reset, chalk.Blue, ss.SpotifyTrack.URI, chalk.Reset)
		fmt.Println()
	}
}
