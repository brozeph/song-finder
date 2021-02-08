// Progam that reads image files from a specified file path,
// then uses Google's ML cloud API for reading text, and
// finally queries Spotify to find matches and create a
// playlist
package main

import (
	"flag"
	"fmt"

	finder "github.com/brozeph/song-finder/internal"
	"github.com/ttacon/chalk"
)

func main() {
	var (
		imageFilePath string
		playlistName  string
	)

	flag.StringVar(&imageFilePath, "path", "", "Specify path to image files.")
	flag.StringVar(&playlistName, "playlist", "", "Provide name for Spotify playlist to create")
	flag.Parse()

	// find all of the image files
	processedFiles, err := finder.Begin(imageFilePath)

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
