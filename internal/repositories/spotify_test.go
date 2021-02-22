package repositories_test

import (
	"testing"

	"github.com/brozeph/song-finder/internal/repositories"
)

func TestSearch(t *testing.T) {
	spotifyRepository := repositories.NewSpotifyRepository()
	if _, err := spotifyRepository.Search("Beck Mixed Business"); err != nil {
		t.Error(err)
	}
}
