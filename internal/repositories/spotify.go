package repositories

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	mrand "math/rand"

	"github.com/brozeph/song-finder/internal/interfaces"
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

type spotifyRepository struct {
	auth          spotify.Authenticator
	clientCH      chan *spotify.Client
	client        *spotify.Client
	codeChallenge string
	codeVerifier  string
	state         string
}

// NewSpotifyRepository returns a new instance
func NewSpotifyRepository() interfaces.ISpotifyRepository {
	return &spotifyRepository{
		auth:     spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistModifyPrivate),
		clientCH: make(chan *spotify.Client),
	}
}

func (r *spotifyRepository) CreatePlaylist(user string, name string, tracks []spotify.SimpleTrack) error {
	_, err := r.ensureClient()
	if err != nil {
		return err
	}

	if _, err := r.client.CreatePlaylistForUser(
		user,
		name,
		"Playlist created by song-finder using image detection of screenshots",
		false); err != nil {
		return err
	}

	return nil
}

func (r *spotifyRepository) Search(searchTerm string) (spotify.SimpleTrack, error) {
	if len(searchTerm) == 0 {
		return spotify.SimpleTrack{}, nil
	}

	_, err := r.ensureClient()
	if err != nil {
		return spotify.SimpleTrack{}, err
	}

	log.Debug().Str("song", searchTerm).Msg("searching for song")

	results, err := r.client.Search(searchTerm, spotify.SearchTypeTrack)
	if err != nil {
		log.Debug().Str("song", searchTerm).Stack().Err(err).Msg("error searching for song")
		return spotify.SimpleTrack{}, err
	}

	if results.Tracks == nil || results.Tracks.Total == 0 {
		log.Debug().Str("song", searchTerm).Msg("no matches found for song")
		return spotify.SimpleTrack{}, nil
	}

	log.Debug().
		Str("song", searchTerm).
		Int("matches", results.Tracks.Total).
		Msg("match(es) found while searching for song")

	return results.Tracks.Tracks[0].SimpleTrack, nil
}

func (r *spotifyRepository) completeAuth(w http.ResponseWriter, res *http.Request) {
	tok, err := r.auth.TokenWithOpts(
		r.state,
		res,
		oauth2.SetAuthURLParam("code_verifier", r.codeVerifier))
	if err != nil {
		http.Error(w, "couldn't get token", http.StatusForbidden)
		log.Debug().Err(err).Msg("couldn't get token")
		panic(err)
	}

	if st := res.FormValue("state"); st != r.state {
		http.NotFound(w, res)
		log.Debug().
			Err(err).
			Str("state", r.state).
			Str("actual state", st).
			Msg("state mismatch")
	}

	cl := r.auth.NewClient(tok)
	fmt.Fprint(w, `<!DOCTYPE html><html lang="en"><head><title>Song Finder: Spotify Auth</title></head><body>
	<p>
		<label>Login process completed</label>
		<div>You may close this window now.</div>
		<div><img src="/assets/img" /></div>
	</p>
</body></html>`)
	r.clientCH <- &cl
}

func (r *spotifyRepository) ensureClient() (*spotify.Client, error) {
	if r.client != nil {
		return r.client, nil
	}

	// ensurer state, codeChallenge and codeVerifier are set
	log.Debug().Msg("setting OAuth params")
	if err := r.setOauthParams(); err != nil {
		return nil, err
	}

	swg := &sync.WaitGroup{}
	swg.Add(1)

	srv := r.startServer(swg)
	log.Debug().Msg("http server started for Spotify authentication flow")

	url := r.auth.AuthURLWithOpts(
		r.state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", r.codeChallenge),
	)

	log.Debug().Str("URL", url).Msg("Spotfy login URL created")
	browser.OpenURL(url)

	cl := <-r.clientCH
	r.client = cl

	// delay a sec to finish serving request for image
	time.Sleep(1 * time.Second)
	if err := srv.Shutdown(context.TODO()); err != nil {
		return r.client, err
	}

	// wait for goroutine in startServer to complete
	swg.Wait()

	user, err := cl.CurrentUser()
	if err != nil {
		return r.client, err
	}

	log.Debug().Str("User.ID", user.ID).Msg("user authenticated")
	return r.client, nil
}

func (r *spotifyRepository) setOauthParams() error {
	// create codeVerifier
	cv, err := randomBytes(
		mrand.Intn(codeVerifierMaxLength-codeVerifierMinLength) + codeVerifierMinLength)
	if err != nil {
		return err
	}

	r.codeVerifier = encode(cv)

	// create codeChallenge from verifier
	h := sha256.New()
	h.Write([]byte(r.codeVerifier))
	r.codeChallenge = encode(h.Sum(nil))

	// create state
	s, err := randomBytes(stateLength)
	if err != nil {
		return err
	}

	r.state = encode(s)

	return nil
}

func (r *spotifyRepository) startServer(wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/callback", r.completeAuth)
	http.HandleFunc("/assets/img", func(w http.ResponseWriter, res *http.Request) {
		i, err := ioutil.ReadFile("assets/b99window.gif")
		if err != nil {
			log.Error().Stack().Err(err).Msg("unable to load funny image")
		}

		w.Header().Set("Content-Type", "image/gif")
		w.Header().Set("Content-Length", strconv.Itoa(len(i)))
		if _, err := w.Write(i); err != nil {
			log.Error().Stack().Err(err).Msg("unable to render funny image")
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, res *http.Request) {
		log.Debug().
			Str("URL", res.URL.String()).
			Msg("received request")
	})

	go func() {
		defer wg.Done()

		if err := srv.ListenAndServe(); err != nil {
			log.Trace().
				Msg("completed Spotify PKCE auth callback")
		}
	}()

	return srv
}

func encode(msg []byte) string {
	encoded := base64.StdEncoding.EncodeToString(msg)
	encoded = strings.Replace(encoded, "+", "-", -1)
	encoded = strings.Replace(encoded, "/", "_", -1)
	encoded = strings.Replace(encoded, "=", "", -1)
	return encoded
}

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
