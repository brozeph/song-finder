package services

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/brozeph/song-finder/internal/interfaces"
	"github.com/brozeph/song-finder/internal/models"

	"github.com/rs/zerolog/log"
	"github.com/superhawk610/bar"
	"github.com/ttacon/chalk"
)

var (
	cr  = regexp.MustCompile(`\n`)
	div = regexp.MustCompile(`(\b|\.)( - )(\b)`)
	dot = regexp.MustCompile(`\s?•\s`)
	//exs     = regexp.MustCompile(`\s{2,}`)
	feat   = regexp.MustCompile(`(?i)\(feat\. [.\s\d\w]*\)`)
	jnk    = regexp.MustCompile(`([•.]\s?){3}`)
	num    = regexp.MustCompile(`^[\d\W]*$`)
	pndra  = regexp.MustCompile(`(?i)\bpandora\b`)
	ply    = regexp.MustCompile(`(?i)^playing from`)
	prp    = regexp.MustCompile(`(?i)po(r)?(n)?tland radi[so] pr(o)?[jy]e[ac]t`)
	rm     = regexp.MustCompile(`(?i)(\w* \S*)?room( \+ [0-9])?$`)
	shzm   = regexp.MustCompile(`(?i)[0-9,]*\s*shazams`)
	sns    = regexp.MustCompile(`(?i)sonos`)
	snsrad = regexp.MustCompile(`(?i)on sonos radio`)
	sp     = regexp.MustCompile(` `)
	sptfy  = regexp.MustCompile(`(?i)[\n\s]spotify\b`)
	swpup  = regexp.MustCompile(`(?i)swipe up to [oó]pen`)
	wd     = regexp.MustCompile(`\w+`)
)

type screenshotService struct {
	screenshotRepository *interfaces.IScreenshotRepository
	spotifyRepository    *interfaces.ISpotifyRepository
	stateRepository      *interfaces.IStateRepository
}

// NewScreenshotService returns new instance of an IScreenshotService
func NewScreenshotService(
	ssr *interfaces.IScreenshotRepository,
	spr *interfaces.ISpotifyRepository,
	str *interfaces.IStateRepository) interfaces.IScreenshotService {

	return &screenshotService{
		screenshotRepository: ssr,
		spotifyRepository:    spr,
		stateRepository:      str,
	}
}

// Begin starts processing the supplied path
// and reading image files
func (ss *screenshotService) Begin(path string) (models.State, error) {
	var (
		ssr   = *ss.screenshotRepository
		spr   = *ss.spotifyRepository
		state models.State
		str   = *ss.stateRepository
	)

	if err := str.Load(state); err != nil {
		if !os.IsNotExist(err) {
			log.Fatal().Stack().Err(err).Msg("unable to load state")
		}

		state = models.State{
			Screenshots:     map[string]*models.Screenshot{},
			SoftwareVersion: softwareVersion,
		}
	}

	// load screenshot paths from screenshotRepository
	screenShots, err := ssr.FindInPath(path)
	if err != nil {
		return state, err
	}

	b := bar.NewWithOpts(
		bar.WithDimensions(len(screenShots), len(screenShots)),
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

	for _, s := range screenShots {
		b.Tick()

		text, err := ssr.DetectText(s.Path)
		if err != nil {
			return state, err
		}

		song := ss.SearchTerm(text)
		track, err := spr.Search(song)

		if err != nil {
			return state, err
		}

		s.LastSearched = time.Now()
		s.SongSearchTerm = song
		s.SpotifyTrack = track

		state.Screenshots[s.SHASum] = s
	}

	b.Done()

	return state, nil
}

// SearchTerm returns a possible artist and
// and song title match from the annotation
func (ss *screenshotService) SearchTerm(annotation string) string {
	var (
		lines     []string
		songParts []string
	)

	isPandora := pndra.MatchString(annotation)
	isPRP := prp.MatchString(annotation)
	isShazam := shzm.MatchString(annotation)
	isSonosRadio := snsrad.MatchString(annotation)
	isSpotify := !isShazam && sptfy.MatchString(annotation)

	// Can safely remove all lines above the first occurrence
	// of Portland Radio Project in the annotation given
	// the location of the song title and artist name in
	// the screen capture
	if isPRP {
		loc := prp.FindStringIndex(annotation)

		if len(loc) > 1 {
			lines = cr.Split(annotation[loc[1]:], -1)
		}
	}

	if lines == nil {
		lines = cr.Split(annotation, -1)
	}

	if len(lines) == 1 {
		return lines[0]
	}

	for i, line := range lines {
		// filter out blank lines
		if line == "" {
			continue
		}

		// when Shazam or Sonos radio, continue until near song artist and name
		if isShazam || isSonosRadio {
			if shzm.MatchString(line) || snsrad.MatchString(line) {
				artist := lines[i-1]
				name := lines[i-2]

				// Shazam wraps multiple artists
				if name[len(name)-1:] == "&" {
					artist = fmt.Sprintf("%s %s", name, artist)
					name = lines[i-3]
				}

				// Sonos radio has the 3 dots
				if jnk.MatchString(name) {
					name = lines[i-3]
				}

				return sanitizeSong(fmt.Sprintf("%s %s", artist, name))
			}

			continue
		}

		// when Spotify or Pandora, continue until near the song artist and name
		if isSpotify || isPandora {
			// check to see if two numbers appear on the same line (scrubber)
			// and that the song is playing from Spotify
			if num.MatchString(line) && num.MatchString(lines[i+1]) {
				artist := lines[i+3]
				name := lines[i+2]

				// handle scenarios where the 3 dots is detected in the image
				if artist == `` || jnk.MatchString(artist) {
					artist = lines[i+4]
				}

				// safe to clear everything prior to this point because the
				// song detail begins below (in fact, the song name is next)
				return formatSongFromSpotifyOrPandora(artist, name)
			}

			continue
		}

		// if the song divider is present on this line,
		// return directly
		if div.MatchString(line) {
			return sanitizeSong(div.ReplaceAllString(line, "$1 $3"))
		}

		// filter out Numbers only
		if num.MatchString(line) {
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

	// join the artist and song name for search
	return sanitizeSong(strings.Join(songParts, " "))
}

func formatSongFromSpotifyOrPandora(artist string, name string) string {
	loc := dot.FindStringIndex(artist)

	if len(loc) > 1 {
		artist = artist[:loc[0]]
	}

	return sanitizeSong(fmt.Sprintf("%s %s", artist, name))
}

func sanitizeSong(song string) string {
	// convert to lowercase
	song = strings.ToLower(song)

	// remove multiple spaces
	song = strings.ReplaceAll(song, "  ", " ")

	// swap & with ,
	song = strings.ReplaceAll(song, " &", ",")

	// remove (feat. XXXX) wording
	song = feat.ReplaceAllString(song, "")

	// removing leading and trailing space
	song = strings.TrimSpace(song)

	// remove the divider char (-) if found
	song = div.ReplaceAllString(song, "")

	return song
}
