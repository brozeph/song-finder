package services_test

import (
	"fmt"
	"testing"

	"github.com/brozeph/song-finder/services"
)

func TestSongFromAnnotationFromPRP(t *testing.T) {
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
	song := services.SongFromAnnotation(testAnnotation)

	if song != "Reverend Freakchild - Personal Jesus (On the..." {
		fmt.Println(song)
		t.Errorf("expected song result from annotation")
	}
}

func TestSongFromAnnotationFromPRPLockScreen(t *testing.T) {
	testAnnotation := `
Portland Radio Project
Pontland Radio Prjeat
The Dig - Soul of the Night

`
	song := services.SongFromAnnotation(testAnnotation)

	if song != "The Dig - Soul of the Night" {
		fmt.Println(song)
		t.Errorf("expected song result from annotation")
	}
}

func TestSongFromAnnotationFromPRPDifferentRooms(t *testing.T) {
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
	song := services.SongFromAnnotation(testAnnotation)

	if song != "Blisses B - Twin Geeks" {
		fmt.Println(song)
		t.Errorf("expected song result from annotation")
	}
}

func TestSongFromAnnotationFromSonosRadio(t *testing.T) {
	testAnnotation := `
Guest Room
Jessy Lanza
Sunset Fuzz on SONOS Radio
Мy Sonos
Pull my hair back

`
	expected := "Jessy Lanza - Pull my hair back"
	song := services.SongFromAnnotation(testAnnotation)

	if song != expected {
		fmt.Println(song)
		t.Errorf("expected song result (%s) from annotation was not matched: %s", expected, song)
	}
}

func TestSongFromAnnotationFromSpotify(t *testing.T) {
	testAnnotation := `
SG Lewis
SG Lewis • Chemicals
Playing from E Spotify
`
	expected := "SG Lewis - Chemicals"
	song := services.SongFromAnnotation(testAnnotation)

	if song != expected {
		fmt.Println(song)
		t.Errorf("expected song result (%s) from annotation was not matched: %s", expected, song)
	}
}

func TestSongFromAnnotationFromSpotifyLockScreen(t *testing.T) {
	testAnnotation := `
Tuesday, February 12
Portland Radio Project
Pnthnt Ruda Pge Smallpools - Stumblin' Home
Swipe up to ópen
`
	expected := "Pnthnt Ruda Pge Smallpools - Stumblin' Home"
	song := services.SongFromAnnotation(testAnnotation)

	if song != expected {
		fmt.Println(song)
		t.Errorf("expected song result (%s) from annotation was not matched: %s", expected, song)
	}
}
