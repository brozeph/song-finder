# Song Finder

## Prerequisites

* Google Cloud API access credentials
* Spotify API access credentials
* Golang

### Setup Google Cloud API Account

See the following: <https://cloud.google.com/sdk/docs/authorizing>

For my purposes, I set an environment variable named `GOOGLE_APPLICATION_CREDENTIALS` pointing to a JSON file with a private key pair for a service account.

### Setup Spotify API Account

See the following: https://developer.spotify.com/dashboard/login

## Execute

```bash
git clone https://github.com/brozeph/song-finder.git
cd song-finder
go run . --path /path/to/images
```

## Testing

```bash
go test ./...
```