// Package services provides utilities for analyzing
// images of songs for the purposes of finding the
// Spotify URI
package services

import (
	"context"
	"os"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

var client spotify.Client

func ensureClient() error {
	if (spotify.Client{} != client) {
		return nil
	}

	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		TokenURL:     spotify.TokenURL,
	}

	token, err := config.Token(context.Background())
	if err != nil {
		return err
	}

	client = spotify.Authenticator{}.NewClient(token)

	return nil
}

// Search looks up a track within the Spotify API and returns
// the first matching result (presumably the most relevant) if
// a match is found
func Search(song string) (spotify.SimpleTrack, error) {
	if err := ensureClient(); err != nil {
		return spotify.SimpleTrack{}, err
	}

	if len(song) == 0 {
		return spotify.SimpleTrack{}, nil
	}

	results, err := client.Search(song, spotify.SearchTypeTrack)
	if err != nil {
		return spotify.SimpleTrack{}, err
	}

	if results.Tracks == nil || results.Tracks.Total == 0 {
		return spotify.SimpleTrack{}, nil
	}

	return results.Tracks.Tracks[0].SimpleTrack, nil
}
