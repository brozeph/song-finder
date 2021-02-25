package repositories

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/brozeph/song-finder/internal/interfaces"
)

var (
	defaultMarshaller = func(v interface{}) (io.Reader, error) {
		b, err := json.MarshalIndent(v, "", "\t")
		if err != nil {
			return nil, err
		}

		return bytes.NewReader(b), nil
	}
	defaultUnmarshaller = func(rdr io.Reader, v interface{}) error {
		return json.NewDecoder(rdr).Decode(v)
	}
)

type stateRepository struct {
	lock      sync.Mutex
	Marshal   func(v interface{}) (io.Reader, error)
	Path      string
	Unmarshal func(rdr io.Reader, v interface{}) error
}

// NewStateRepository returns an instance of IStateRepository for
// saving and loading application state
func NewStateRepository(path string) interfaces.IStateRepository {
	return &stateRepository{
		Marshal:   defaultMarshaller,
		Path:      path,
		Unmarshal: defaultUnmarshaller,
	}
}

// Load retrieves a persisted state for use
func (r *stateRepository) Load(v interface{}) error {
	// lock for thread safety
	r.lock.Lock()
	defer r.lock.Unlock()

	// open the state file
	f, err := os.Open(r.Path)
	if err != nil {
		return err
	}

	defer f.Close()

	// return the unmarshalled object
	return r.Unmarshal(f, v)
}

// Save persists state for a subsequent run
func (r *stateRepository) Save(v interface{}) error {
	// lock for thread safety
	r.lock.Lock()
	defer r.lock.Unlock()

	// create the file at specified path
	fil, err := os.Create(r.Path)
	if err != nil {
		return err
	}

	defer fil.Close()

	// marshal the object
	rdr, err := r.Marshal(v)
	if err != nil {
		return err
	}

	// write to file
	_, err = io.Copy(fil, rdr)

	return err
}
