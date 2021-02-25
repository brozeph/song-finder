package interfaces

import "github.com/brozeph/song-finder/internal/models"

// IScreenshotService provides the workflow for processing screenshots
// and creating Spotify playlists
type IScreenshotService interface {
	Begin(path string) (models.State, error)
	SearchTerm(annotation string) string
}
