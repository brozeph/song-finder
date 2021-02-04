// Package services provides utilities for analyzing
// images of songs for the purposes of finding the
// Spotify URI
package services

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	vision "cloud.google.com/go/vision/apiv1"
)

var (
	cr    *regexp.Regexp = regexp.MustCompile(`\n`)
	div   *regexp.Regexp = regexp.MustCompile(`(\w)(\ \-\ ){1}(\w)`)
	dot   *regexp.Regexp = regexp.MustCompile(`\s?•\s`)
	exs   *regexp.Regexp = regexp.MustCompile(`\s{2,}`)
	num   *regexp.Regexp = regexp.MustCompile(`^[\d\W]*$`)
	ply   *regexp.Regexp = regexp.MustCompile(`(?i)^playing from`)
	prp   *regexp.Regexp = regexp.MustCompile(`(?i)po(r)?(n)?tland radi[so] pr(o)?[jy]e[ac]t`)
	rm    *regexp.Regexp = regexp.MustCompile(`(?i)(\w*(\ ){1}\S*){0,1}room(\ \+\ [0-9]){0,1}$`)
	sns   *regexp.Regexp = regexp.MustCompile(`(?i)sonos`)
	sp    *regexp.Regexp = regexp.MustCompile(` `)
	sptfy *regexp.Regexp = regexp.MustCompile(`\nSpotify\n`)
	swpup *regexp.Regexp = regexp.MustCompile(`(?i)swipe up to [oó]pen`)
	wd    *regexp.Regexp = regexp.MustCompile(`\w+`)
)

func dedup(value string) string {
	found := map[string]bool{}
	words := []string{}

	for _, word := range sp.Split(value, -1) {
		word = strings.ToLower(word)

		if found[word] {
			continue
		}

		found[word] = true
		words = append(words, word)
	}

	return strings.Join(words, " ")
}

func formatSongFromSpotify(artist string, name string) string {
	loc := dot.FindStringIndex(artist)

	fmt.Println(artist)
	fmt.Println(name)

	if len(loc) > 1 {
		artist = artist[:loc[0]]
	}

	return strings.ToLower(strings.Join([]string{artist, name}, " "))
}

func stripRunes(value string) string {
	stripped := []rune{}

	for _, ch := range value {
		if utf8.RuneLen(ch) > 2 {
			continue
		}

		stripped = append(stripped, ch)
	}

	return string(stripped)
}

// DetectText accepts an image path, reads the image and
// requests to retrieve text annotations from the Google Cloud
// vision API
func DetectText(imagePath string) (string, error) {
	text := ""
	ctx := context.Background()

	f, err := os.Open(imagePath)
	if err != nil {
		return text, err
	}

	ifr, err := vision.NewImageFromReader(f)
	if err != nil {
		return text, err
	}

	// for each image, upload to Google image analysis
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return text, err
	}

	annotations, err := client.DetectTexts(ctx, ifr, nil, 10)
	if err != nil {
		return text, err
	}

	if annotations[0] != nil {
		text = annotations[0].Description
	}

	fmt.Println(text)

	return text, nil
}

// SongArtistAndName returns a possible artist and
// and song title match from the annotation
func SongArtistAndName(annotation string) string {
	var (
		lines         []string
		songParts     []string
		subAnnotation string
	)

	isSpotify := sptfy.MatchString(annotation)

	// Can safely remove all lines above the first occurrence
	// of Portland Radio Project in the annotation given
	// the location of the song title and artist name in
	// the screen capture
	loc := prp.FindStringIndex(annotation)

	if len(loc) > 1 {
		subAnnotation = annotation[loc[1]:]
	} else {
		subAnnotation = annotation
	}

	lines = cr.Split(subAnnotation, -1)

	if len(lines) == 1 {
		return lines[0]
	}

	for i, line := range lines {
		// filter out blank lines
		if line == "" {
			continue
		}

		// if the song divider is present on this line,
		// return directly
		if div.MatchString(line) {
			return strings.ToLower(div.ReplaceAllString(line, "$1 $3"))
		}

		// filter out Numbers only
		if num.MatchString(line) {
			// check to see if two numbers appear on the same line (scrubber)
			// and that the song is playing from Spotify
			if isSpotify && num.MatchString(lines[i+1]) {
				// safe to clear everything prior to this point because the
				// song detail begins below (in fact, the song name is next)
				return formatSongFromSpotify(lines[i+3], lines[i+2])
			}

			continue
		}

		// filter out "playing from ..."
		if ply.MatchString(line) {
			continue
		}

		// filter out Portland Radio Project
		if prp.MatchString(line) {
			continue
		}

		// filter out Sonos room labels
		if rm.MatchString(line) {
			continue
		}

		// filiter out lines w/ Sonos
		if sns.MatchString(line) {
			continue
		}

		// filter out lines w/o spaces
		if !sp.MatchString(line) {
			continue
		}

		// filter out lines that say "swipe up to open" (iPhone)
		if swpup.MatchString(line) {
			continue
		}

		// filter out non-word lines
		if !wd.MatchString(line) {
			continue
		}

		songParts = append(songParts, line)
	}

	song := strings.Join(songParts, " ")
	song = stripRunes(song)
	song = dedup(song)
	song = strings.ReplaceAll(song, "  ", " ")

	return song
}
