// Package services provides utilities for analyzing
// images of songs for the purposes of finding the
// Spotify URI
package services

import (
	"regexp"
	"strings"
)

var (
	cr    *regexp.Regexp = regexp.MustCompile(`\n`)
	num   *regexp.Regexp = regexp.MustCompile(`^[\d\W]*$`)
	prp   *regexp.Regexp = regexp.MustCompile(`(?i)po(r)?(n)?tland radi[so] pr(o)?[jy]e[ac]t`)
	rm    *regexp.Regexp = regexp.MustCompile(`(?i)(\w*(\ ){1}\S*){0,1}room(\ \+\ [0-9]){0,1}$`)
	sonos *regexp.Regexp = regexp.MustCompile(`(?i)sonos`)
	sp    *regexp.Regexp = regexp.MustCompile(` `)
	wd    *regexp.Regexp = regexp.MustCompile(`\w+`)
)

/*
func dedup(value string) string {
	words := []string

	for _, word range sp.Split(value, -1) {

	}

	return ""
}
*/

// SongFromAnnotation returns a possible artist and
// and song title match from the annotation
func SongFromAnnotation(annotation string) string {
	var (
		lines         []string
		song          []string
		subAnnotation string
	)

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

	for _, line := range lines {
		// filter out blank lines
		if line == "" {
			continue
		}

		// filter out Numbers only
		if num.MatchString(line) {
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
		if sonos.MatchString(line) {
			continue
		}

		// filter out lines w/o spaces
		if !sp.MatchString(line) {
			continue
		}

		// filter out non-word lines
		if !wd.MatchString(line) {
			continue
		}

		song = append(song, line)
	}

	return strings.Join(song, " - ")
}
