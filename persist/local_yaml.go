package persist

import (
	"errors"
	"io"
	"os"

	"github.com/gofrs/flock"
	"github.com/kiyutink/tipam/tipam"
	"gopkg.in/yaml.v3"
)

const (
	defaultYAMLStateAPIVersion = 1
)

type yamlStateClaim struct {
	Tags  []string `yaml:"tags"`
	Final bool     `yaml:"final,omitempty"`
}

type yamlState struct {
	APIVersion int                       `yaml:"apiVersion"`
	Claims     map[string]yamlStateClaim `yaml:"claims,omitempty"`
}

func newEmptyYAMLState() *yamlState {
	return &yamlState{
		APIVersion: defaultYAMLStateAPIVersion,
		Claims:     map[string]yamlStateClaim{},
	}
}

type LocalYAML struct {
	fileName string
	flock    *flock.Flock
}

func NewLocalYAML(fileName string) *LocalYAML {
	return &LocalYAML{
		fileName: fileName,
		flock:    flock.New(fileName),
	}
}

func (lyp *LocalYAML) Persist(s *tipam.State) error {
	file, err := os.OpenFile(lyp.fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	ys := stateToYAMLState(s)

	err = encodeYAMLState(ys, file)
	return err
}

func (lyp *LocalYAML) Read() (*tipam.State, error) {
	bytes, err := os.ReadFile(lyp.fileName)

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

func (lyp *LocalYAML) Lock() error {
	return lyp.flock.Lock()
}

func (lyp *LocalYAML) Unlock() error {
	return lyp.flock.Unlock()
}

func yamlStateToState(ys *yamlState) (*tipam.State, error) {
	s := tipam.NewState()

	for cidr, yc := range ys.Claims {
		c, err := tipam.ParseClaimFromCIDR(cidr, yc.Tags, yc.Final)
		if err != nil {
			return nil, err
		}
		s.Claims[cidr] = c
	}

	return s, nil
}

func stateToYAMLState(s *tipam.State) *yamlState {
	ys := newEmptyYAMLState()
	for cidr, r := range s.Claims {
		ys.Claims[cidr] = yamlStateClaim{Tags: r.Tags, Final: r.Final}
	}
	return ys
}

func encodeYAMLState(ys *yamlState, w io.Writer) error {
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)

	return encoder.Encode(ys)
}

func decodeYAMLState(ys *yamlState, r io.Reader) error {
	decoder := yaml.NewDecoder(r)

	return decoder.Decode(ys)
}
