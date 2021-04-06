package interfaces

import (
	"github.com/brozeph/song-finder/internal/models"
	"github.com/zmb3/spotify"
)

type IPlaylistService interface {
	EnsurePlaylist(name string, tracks []spotify.SimpleTrack) error
}

// IScreenshotService provides the workflow for processing screenshots
// and creating Spotify playlists
type IScreenshotService interface {
	Begin(path string) (models.State, error)
	SearchTerm(annotation string) string
}
