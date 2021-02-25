package models

import (
	"time"

	"github.com/zmb3/spotify"
)

// Screenshot contains the details / state for every
// screenshot image being processed
type Screenshot struct {
	LastSearched   time.Time
	Path           string
	SongSearchTerm string
	SpotifyTrack   spotify.SimpleTrack
}
