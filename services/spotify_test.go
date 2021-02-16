package services_test

import (
	"testing"

	"github.com/brozeph/song-finder/services"
)

func TestSearch(t *testing.T) {
	if _, err := services.Search("Beck Mixed Business"); err != nil {
		t.Error(err)
	}
}
