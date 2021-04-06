package services

import (
	"github.com/brozeph/song-finder/internal/interfaces"
	"github.com/zmb3/spotify"
)

type playlistService struct {
	spotifyRepository *interfaces.ISpotifyRepository
	stateRepository   *interfaces.IStateRepository
}

func NewPlaylistService(
	spr *interfaces.ISpotifyRepository,
	str *interfaces.IStateRepository) interfaces.IPlaylistService {
	return playlistService{
		spotifyRepository: spr,
		stateRepository:   str,
	}
}

func (ps playlistService) EnsurePlaylist(name string, tracks []spotify.SimpleTrack) error {
	return nil
}

func (ps playlistService) lookupPlaylist(name string) spotify.SimplePlaylist {
	return spotify.SimplePlaylist{}
}
