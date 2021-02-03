package services_test

import (
	"testing"

	"github.com/brozeph/song-finder/services"
)

func TestSongArtistAndNameFromPRP(t *testing.T) {
	testAnnotation := `
11:53
Master Bedroom + 3
BP
Porntland Radis Pryeat
Portland Radio Project
Reverend Freakchild - Personal Jesus (On the...
Portland Radio Project on TUNE
IN
K
») -
☆
Мy Sonos
Browse
Rooms
Search
Settings

`
	expected := "reverend freakchild personal jesus (on the..."
	song := services.SongArtistAndName(testAnnotation)

	if song != expected {
		t.Errorf("expected song result (%s) from annotation was not matched: %s", expected, song)
	}
}

func TestSongArtistAndNameFromPRPLockScreen(t *testing.T) {
	testAnnotation := `
Portland Radio Project
Pontland Radio Prjeat
The Dig - Soul of the Night

`
	expected := "the dig soul of the night"
	song := services.SongArtistAndName(testAnnotation)

	if song != expected {
		t.Errorf("expected song result (%s) from annotation was not matched: %s", expected, song)
	}
}

func TestSongArtistAndNameFromPRPDifferentRooms(t *testing.T) {
	testAnnotation := `
3:40 1
Move + Den
CD 100%
BP
Porntland Radis Pryeat
Portland Radio Project
• . .
Blisses B - Twin Geeks
Portland Radio Project on TUNE
IN
K
>)
☆
Мy Sonos
Browse
Rooms
Search
Settings

`

	expected := "blisses b twin geeks"
	song := services.SongArtistAndName(testAnnotation)

	if song != expected {
		t.Errorf("expected song result (%s) from annotation was not matched: %s", expected, song)
	}
}

func TestSongArtistAndNameFromSonosRadio(t *testing.T) {
	testAnnotation := `
Guest Room
Jessy Lanza
Sunset Fuzz on SONOS Radio
Мy Sonos
Pull my hair back

`
	expected := "jessy lanza pull my hair back"
	song := services.SongArtistAndName(testAnnotation)

	if song != expected {
		t.Errorf("expected song result (%s) from annotation was not matched: %s", expected, song)
	}
}

func TestSongArtistAndNameFromSpotify(t *testing.T) {
	testAnnotation := `
SG Lewis
SG Lewis • Chemicals
Playing from E Spotify
`
	expected := "sg lewis chemicals"
	song := services.SongArtistAndName(testAnnotation)

	if song != expected {
		t.Errorf("expected song result (%s) from annotation was not matched: %s", expected, song)
	}
}

func TestSongArtistAndNameFromSpotifyLockScreen(t *testing.T) {
	testAnnotation := `
Tuesday, February 12
Portland Radio Project
Pnthnt Ruda Pge Smallpools - Stumblin' Home
Swipe up to ópen
`
	expected := "pnthnt ruda pge smallpools stumblin' home"
	song := services.SongArtistAndName(testAnnotation)

	if song != expected {
		t.Errorf("expected song result (%s) from annotation was not matched: %s", expected, song)
	}
}

func TestSongArtistAndNameFromSpotifyInSonos(t *testing.T) {
	testAnnotation := `
6:46 1
Kitchen + 2
ZD 100%
SONNY ALVEN
WASTĘD YOUTH (FEAT. CAL)
AMERIC
1:11
-2:05
Wasted Youth
Sonny Alven• Girls - EP
feel good
Spotify
!!
%3D

`
	expected := "Sonny Alven Girls EP"
	song := services.SongArtistAndName(testAnnotation)

	if song != expected {
		t.Errorf("expected song result (%s) from annotation was not matched: %s", expected, song)
	}
}
