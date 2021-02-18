// Package services provides utilities for analyzing
// images of songs for the purposes of finding the
// Spotify URI
package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	mrand "math/rand"
	"net/http"
	"strings"
	"sync"

	"github.com/pkg/browser"
	"github.com/rs/zerolog/log"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

const (
	codeVerifierMaxLength = 128
	codeVerifierMinLength = 43
	redirectURI           = "http://localhost:8080/callback"
	stateLength           = 36
)

var (
	auth          = spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistModifyPrivate)
	cch           = make(chan *spotify.Client)
	client        *spotify.Client
	state         string
	codeChallenge string
	codeVerifier  string
)

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.TokenWithOpts(
		state,
		r,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		http.Error(w, "couldn't get token", http.StatusForbidden)
		log.Debug().Err(err).Msg("couldn't get token")
		panic(err)
	}

	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Debug().Err(err).Str("state", state).Str("actual state", st).Msg("state mismatch")
	}

	cl := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	cch <- &cl
}

func encode(msg []byte) string {
	encoded := base64.StdEncoding.EncodeToString(msg)
	encoded = strings.Replace(encoded, "+", "-", -1)
	encoded = strings.Replace(encoded, "/", "_", -1)
	encoded = strings.Replace(encoded, "=", "", -1)
	return encoded
}

func ensureClient() (*spotify.Client, error) {
	if client != nil {
		return client, nil
	}

	// ensurer state, codeChallenge and codeVerifier are set
	log.Debug().Msg("setting OAuth params")
	if err := setOauthParams(); err != nil {
		return nil, err
	}

	swg := &sync.WaitGroup{}
	swg.Add(1)

	srv := startServer(swg)
	log.Debug().Msg("http server started for Spotify authentication flow")

	url := auth.AuthURLWithOpts(
		state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	)

	log.Debug().Str("URL", url).Msg("Spotfy login URL created")
	browser.OpenURL(url)

	cl := <-cch
	client = cl

	if err := srv.Shutdown(context.TODO()); err != nil {
		return client, err
	}

	// wait for goroutine in startServer to complete
	swg.Wait()

	user, err := cl.CurrentUser()
	if err != nil {
		return client, err
	}

	log.Debug().Str("User.ID", user.ID).Msg("user authenticated")
	return client, nil
}

/*
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
//*/

// https://tools.ietf.org/html/rfc7636#section-4.1)
func randomBytes(length int) ([]byte, error) {
	const charset = ".ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	const csLen = byte(len(charset))
	output := make([]byte, 0, length)
	for {
		buf := make([]byte, length)
		if _, err := io.ReadFull(rand.Reader, buf); err != nil {
			return nil, fmt.Errorf("failed to read random bytes: %v", err)
		}
		for _, b := range buf {
			// Avoid bias by using a value range that's a multiple of 62
			if b < (csLen * 4) {
				output = append(output, charset[b%csLen])

				if len(output) == length {
					return output, nil
				}
			}
		}
	}
}

func setOauthParams() error {
	// create codeVerifier
	cv, err := randomBytes(
		mrand.Intn(codeVerifierMaxLength-codeVerifierMinLength) + codeVerifierMinLength)
	if err != nil {
		return err
	}

	codeVerifier = encode(cv)

	// create codeChallenge from verifier
	h := sha256.New()
	h.Write([]byte(codeVerifier))
	codeChallenge = encode(h.Sum(nil))

	// create state
	s, err := randomBytes(stateLength)
	if err != nil {
		return err
	}

	state = encode(s)

	return nil
}

func startServer(wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Str("URL", r.URL.String()).Msg("received request")
	})

	go func() {
		defer wg.Done()

		if err := srv.ListenAndServe(); err != nil {
			log.Debug().Err(err).Msg("")
			// panic(err)
		}
	}()

	return srv
}

// CreatePlaylist creates a new Spotify playlist with the
// supplied tracks
func CreatePlaylist(user string, name string, tracks []spotify.SimpleTrack) error {
	_, err := ensureClient()
	if err != nil {
		return err
	}

	if _, err := client.CreatePlaylistForUser(
		user,
		name,
		"Playlist created by song-finder using image detection of screenshots",
		false); err != nil {
		return err
	}

	return nil
}

// Search looks up a track within the Spotify API and returns
// the first matching result (presumably the most relevant) if
// a match is found
func Search(song string) (spotify.SimpleTrack, error) {
	_, err := ensureClient()
	if err != nil {
		return spotify.SimpleTrack{}, err
	}

	if len(song) == 0 {
		return spotify.SimpleTrack{}, nil
	}

	log.Debug().Str("song", song).Msg("searching for song")

	results, err := client.Search(song, spotify.SearchTypeTrack)
	if err != nil {
		log.Debug().Str("song", song).Stack().Err(err).Msg("error searching for song")
		return spotify.SimpleTrack{}, err
	}

	if results.Tracks == nil || results.Tracks.Total == 0 {
		log.Debug().Str("song", song).Msg("no matches found for song")
		return spotify.SimpleTrack{}, nil
	}

	log.Debug().
		Str("song", song).
		Int("matches", results.Tracks.Total).
		Msg("match(es) found while searching for song")

	return results.Tracks.Tracks[0].SimpleTrack, nil
}
