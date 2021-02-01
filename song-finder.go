// Progam that reads image files from a specified file path
// and then uses Google's ML cloud API for reading text
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/brozeph/song-finder/services"
)

var imageExtensions = []string{".png", ".jpg", ".jpeg"}

func isImage(filePath string) bool {
	var fileExt = strings.ToLower(filepath.Ext(filePath))

	for _, ext := range imageExtensions {
		if fileExt == ext {
			return true
		}
	}

	return false
}

func main() {
	var (
		imageFiles    []string
		imageFilePath string
	)

	flag.StringVar(&imageFilePath, "path", "", "Specify path to image files.")
	flag.Parse()

	// find all of the image files
	err := filepath.Walk(imageFilePath, func(path string, info os.FileInfo, err error) error {
		if isImage(path) && info.Size() > 0 {
			imageFiles = append(imageFiles, path)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	for _, imageFile := range imageFiles {
		fmt.Println("Processing image at:", imageFile)

		text, err := services.DetectText(imageFile)
		if err != nil {
			panic(err)
		}

		song := services.SongArtistAndName(text)

		fmt.Println(song)
	}
}
