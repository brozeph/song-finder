package interfaces

import "github.com/zmb3/spotify"

// IScreenshotRepository provides methods for retrieving screenshots
// from the filesystem
type IScreenshotRepository interface {
	DetectText(path string) (string, error)
	FindInPath(path string) ([]string, error)
}

// ISpotifyRepository provides methods to abstract interaction with the
// Spotify API
type ISpotifyRepository interface {
	CreatePlaylist(user string, name string, tracks []spotify.SimpleTrack) error
	Search(searchTerm string) (spotify.SimpleTrack, error)
}

// IStateRepository provides methods to persist and retrieve state
// for subsequent runs of the application
type IStateRepository interface {
	Load(v interface{}) error
	Save(v interface{}) error
}
