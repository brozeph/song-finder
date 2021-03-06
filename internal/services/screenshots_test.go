package services_test

import (
	"testing"

	"github.com/brozeph/song-finder/internal/services"
)

var s = services.NewScreenshotService(nil, nil, nil)

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
	song := s.SearchTerm(testAnnotation)

	if song != expected {
		t.Errorf("expected song result \"%s\" from annotation was not matched: \"%s\"", expected, song)
	}
}

func TestSongArtistAndNameFromPRP2(t *testing.T) {
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
	song := s.SearchTerm(testAnnotation)

	if song != expected {
		t.Errorf("expected song result \"%s\" from annotation was not matched: \"%s\"", expected, song)
	}
}

func TestSongArtistAndNameFromPRP3(t *testing.T) {
	testAnnotation := `8:47 1
Kitchen
PP
Pontland Rad9 Projedt
Portland Radio Project
•..
YellowStraps - Goldress (feat. VYNK)
Portland Radio Project
IN
TUNE
K]
!!
`
	expected := "yellowstraps goldress"
	song := s.SearchTerm(testAnnotation)

	if song != expected {
		t.Errorf("expected song result \"%s\" from annotation was not matched: \"%s\"", expected, song)
	}
}

func TestSongArtistAndNameFromPRP4(t *testing.T) {
	testAnnotation := `
11:55
レ
Living Room + 3
ニ
BP
Portland Racfi9 Projedt
Portland Radio Project
RAC Pron R.A.C. - This Song Feat Rostam
Portland Radio Project on TUNE
IN
») -
My Sonos
Browse
Rooms
Search
Settings
`
	expected := "rac pron r.a.c. this song feat rostam"
	song := s.SearchTerm(testAnnotation)

	if song != expected {
		t.Errorf("expected song result \"%s\" from annotation was not matched: \"%s\"", expected, song)
	}
}

func TestSongArtistAndNameFromPRPLockScreen(t *testing.T) {
	testAnnotation := `
Portland Radio Project
Pontland Radio Prjeat
The Dig - Soul of the Night

`
	expected := "the dig soul of the night"
	song := s.SearchTerm(testAnnotation)

	if song != expected {
		t.Errorf("expected song result \"%s\" from annotation was not matched: \"%s\"", expected, song)
	}
}

func TestSongArtistAndNameFromShazam(t *testing.T) {
	testAnnotation := `
SONGWRITER POF
ou
orden
Campe
cemn
a Pen
yougs dhe
In My Atmosphere
Raphael Lake & Eric Brooks &
Camden Rose
2 7,392 Shazams
A Spotify
ОPEN
ADD TO
TOP SONGS
Ready
Raprerel Lake QAaronLovy &...
INDIE SOUL
Lone

`
	expected := "raphael lake, eric brooks, camden rose in my atmosphere"
	song := s.SearchTerm(testAnnotation)

	if song != expected {
		t.Errorf("expected song result \"%s\" from annotation was not matched: \"%s\"", expected, song)
	}
}

func TestSongArtistAndNameFromSonosRadio(t *testing.T) {
	testAnnotation := `
6:31 1
Guest Room
Giddy
• . .
Jessy Lanza
Sunset Fuzz on SONOS Radio
K
>)
☆
Мy Sonos
Browse
Rooms
Search
Settings
JESSY LANZA
Pull my hair back

`
	expected := "jessy lanza giddy"
	song := s.SearchTerm(testAnnotation)

	if song != expected {
		t.Errorf("expected song result \"%s\" from annotation was not matched: \"%s\"", expected, song)
	}
}

func TestSongArtistAndNameFromSonosRadio2(t *testing.T) {
	testAnnotation := `
6:16 1
Guest Room
Obsessed
• . .
Hatchie
Sunset Fuzz on SONOS Radio
K
||
>)
☆
Мy Sonos
Browse
Rooms
Search
Settings
`
	expected := "hatchie obsessed"
	song := s.SearchTerm(testAnnotation)

	if song != expected {
		t.Errorf("expected song result \"%s\" from annotation was not matched: \"%s\"", expected, song)
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
	song := s.SearchTerm(testAnnotation)

	if song != expected {
		t.Errorf("expected song result \"%s\" from annotation was not matched: \"%s\"", expected, song)
	}
}

func TestSongArtistAndNameFromSpotifyMinimal(t *testing.T) {
	testAnnotation := `
SG Lewis
CHEMICALS
1:12
-3:02
Chemicals
• ..
SG Lewis • Chemicals
Playing from E Spotify
`
	expected := "sg lewis chemicals"
	song := s.SearchTerm(testAnnotation)

	if song != expected {
		t.Errorf("expected song result \"%s\" from annotation was not matched: \"%s\"", expected, song)
	}
}

func TestSongArtistAndNameFromSpotify(t *testing.T) {
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
	expected := "sonny alven wasted youth"
	song := s.SearchTerm(testAnnotation)

	if song != expected {
		t.Errorf("expected song result \"%s\" from annotation was not matched: \"%s\"", expected, song)
	}
}

func TestSongArtistAndNameFromSpotify2(t *testing.T) {
	testAnnotation := `
6:06
Guest Room + Den + 4
CD 100%
ZD)
TRUE COLORS
6:55
-0:28
Раpercut
• • •
Zedd • True Colors
Playing from 6 Spotify
K |
») –
☆
Мy Sonos
Browse
Rooms
Search
Settings

`
	expected := "zedd раpercut"
	song := s.SearchTerm(testAnnotation)

	if song != expected {
		t.Errorf("expected song result \"%s\" from annotation was not matched: \"%s\"", expected, song)
	}
}

func TestSongArtistAndNameFromPandora(t *testing.T) {
	testAnnotation := `
8:32 1
Kitchen + 2
ZD 100%
A FINE FRENZY
ONE CELL IN THE SEA
2:22
-1:54
Hope For The Hopeless
A Fine Frenzy • One Cell In The Sea
Rosi Golan Radio (My Station)
pandora
K]
!!
%3D

`
	expected := "a fine frenzy hope for the hopeless"
	song := s.SearchTerm(testAnnotation)

	if song != expected {
		t.Errorf("expected song result \"%s\" from annotation was not matched: \"%s\"", expected, song)
	}
}

func TestSongArtistAndNameFromLinn(t *testing.T) {
	testAnnotation := `
Search
KYDD
Walk On You
NO
WALKING
Кydd
Walk On You
Walk On You
1:16
PCM 44.1 kHz/16 bit 1.4 Mbps
4:30
K ||
44

`
	expected := "кydd walk on you"
	song := s.SearchTerm(testAnnotation)

	if song != expected {
		t.Errorf("expected song result \"%s\" from annotation was not matched: \"%s\"", expected, song)
	}
}

func TestSongArtistAndNameFromLinn2(t *testing.T) {
	testAnnotation := `
Search
DUSTY
LEIGH
NIGHT
PARENTAL
ADVISORY
EXPLICIT CONTENT
Dusty Leigh
No Phucks (feat. Bubba
Sparxxx & Fishscales)
Boujee Nights
0:33
PCM 44.1 kHz/16 bit 1.4 Mbps
3:17
|L
46

`
	expected := "dusty leigh no phucks (feat. bubba sparxxx, fishscales)"
	song := s.SearchTerm(testAnnotation)

	if song != expected {
		t.Errorf("expected song result \"%s\" from annotation was not matched: \"%s\"", expected, song)
	}
}
