// Progam that reads image files from a specified file path,
// then uses Google's ML cloud API for reading text, and
// finally queries Spotify to find matches and create a
// playlist
package main

import (
	"fmt"
	"os"

	finder "github.com/brozeph/song-finder/internal"
	"github.com/brozeph/song-finder/services"
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
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	//*
	song, err := services.Search("Beck Mixed Business")
	if err != nil {
		log.Error().Err(err).Msg("")
		panic(err)
	}
	log.Info().Str("Spotify URI", song.ID.String()).Msg("song found")
	//*/

	var (
		options cmdlineOptions
		parser  = flags.NewParser(&options, flags.Default)
	)

	// parse command line arguments
	if _, err := parser.Parse(); err != nil {
		log.Error().Err(err).Msg("")
		os.Exit(1)
	}

	// find all of the image files
	processedFiles, err := finder.Begin(options.ImageFilePath)

	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Printf(
		"Process completed for %s%d%s files",
		chalk.Blue,
		len(processedFiles),
		chalk.Reset)

	for _, file := range processedFiles {
		fmt.Println(chalk.Blue, "File:", chalk.Reset, file.ImagePath)
		fmt.Println(chalk.Red, "Song:", chalk.Reset, file.SongArtistAndName)
		fmt.Println(chalk.Green, "Spotify URI:", chalk.Reset, chalk.Blue, file.SpotifyTrack.URI, chalk.Reset)
		fmt.Println()
	}
}
