package persist

import (
	"errors"
	"net"
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

	yamlState := newEmptyYAMLState()

	for cidr, r := range s.Claims {
		yamlState.Claims[cidr] = yamlStateClaim{Tags: r.Tags, Final: r.Final}
	}

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)

	err = encoder.Encode(yamlState)

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

	yamlState := newEmptyYAMLState()

	err = yaml.Unmarshal(bytes, yamlState)

	if err != nil {
		return nil, err
	}

	state := &tipam.State{
		Claims: map[string]tipam.Claim{},
	}

	for c, claim := range yamlState.Claims {
		_, ipNet, err := net.ParseCIDR(c)
		if err != nil {
			return nil, err
		}
		r := tipam.NewClaim(ipNet, claim.Tags)
		state.Claims[c] = r
	}

	return state, nil
}

func (lyp *LocalYAML) Lock() error {
	return lyp.flock.Lock()
}

func (lyp *LocalYAML) Unlock() error {
	return lyp.flock.Unlock()
}
