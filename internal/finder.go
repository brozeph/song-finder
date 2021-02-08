// Package finder processes images files and calls
// various services to retrieve details
package finder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/brozeph/song-finder/services"
	"github.com/superhawk610/bar"
	"github.com/ttacon/chalk"
	"github.com/zmb3/spotify"
)

var imageExtensions = []string{".png", ".jpg", ".jpeg"}

// File contains the details / state
// for every image processed
type File struct {
	ImagePath         string
	SongArtistAndName string
	SpotifyTrack      spotify.SimpleTrack
}

// Begin starts processing the supplied path
// and reading image files
func Begin(imageFilePath string) ([]File, error) {
	var (
		imageFiles     []string
		processedFiles []File
	)

	// find all of the image files
	err := filepath.Walk(imageFilePath, func(path string, info os.FileInfo, err error) error {
		if isImage(path) && info.Size() > 0 {
			imageFiles = append(imageFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	b := bar.NewWithOpts(
		bar.WithDimensions(len(imageFiles), len(imageFiles)),
		bar.WithFormat(
			fmt.Sprintf(
				" %sprocessing...%s :percent :bar %s:eta%s     ",
				chalk.Blue,
				chalk.Reset,
				chalk.Green,
				chalk.Reset,
			),
		),
	)

	for _, imageFile := range imageFiles {
		b.Tick()

		text, err := services.DetectText(imageFile)
		if err != nil {
			return nil, err
		}

		song := services.SongArtistAndName(text)
		track, err := services.Search(song)

		if err != nil {
			return nil, err
		}

		processedFiles = append(processedFiles, File{
			ImagePath:         imageFile,
			SongArtistAndName: song,
			SpotifyTrack:      track,
		})
	}

	b.Done()

	return processedFiles, nil
}

func isImage(filePath string) bool {
	var fileExt = strings.ToLower(filepath.Ext(filePath))

	for _, ext := range imageExtensions {
		if fileExt == ext {
			return true
		}
	}

	return false
}
