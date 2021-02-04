# Song Finder

## Prerequisites

* Google Cloud API access credentials
* Spotify API access credentials
* Golang

### Setup Google Cloud API Account

See the following: <https://cloud.google.com/sdk/docs/authorizing>

For my purposes, I set an environment variable named `GOOGLE_APPLICATION_CREDENTIALS` pointing to a JSON file with a private key pair for a service account. The vision API should be enabled for the Google application credentials that are configured.

### Setup Spotify API Account

See the following: https://developer.spotify.com/dashboard/login

The application reads environment variables for `SPOTIFY_CLIENT_ID` and `SPOTIFY_CLIENT_SECRET` and uses these client credentials to authenticate API requests to Spotify.

## Setup

```bash
git clone https://github.com/brozeph/song-finder.git
cd song-finder
go mod download
```

### Running the App

Note the path below should be to the screenshots of captured songs to be looked up.

```bash
go run . --path /path/to/images
```

### Running Tests

```bash
go test ./...
```