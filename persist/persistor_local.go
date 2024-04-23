package persist

import (
	"errors"
	"os"

	"github.com/gofrs/flock"
	"github.com/kiyutink/tipam/tipam"
	"gopkg.in/yaml.v3"
)

type Local struct {
	fileName string
	flock    *flock.Flock
}

func NewLocal(fileName string) *Local {
	return &Local{
		fileName: fileName,
		flock:    flock.New(fileName),
	}
}

func (l *Local) Persist(s *tipam.State) error {
	file, err := os.OpenFile(l.fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	ys := stateToYAMLState(s)

	err = encodeYAMLState(ys, file)
	return err
}

func (l *Local) Read() (*tipam.State, error) {
	bytes, err := os.ReadFile(l.fileName)

	switch {

	case errors.Is(err, os.ErrNotExist):
		fallthrough
	case err == nil:
		// Do nothing
	default:
		return nil, err
	}

	ys := newEmptyYAMLState()

	err = yaml.Unmarshal(bytes, ys)
	if err != nil {
		return nil, err
	}

	s, err := yamlStateToState(ys)

	return s, err
}

func (l *Local) Lock() error {
	return l.flock.Lock()
}

func (l *Local) Unlock() error {
	return l.flock.Unlock()
}
