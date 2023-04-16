package persist

import (
	"errors"
	"net"
	"os"

	"github.com/gofrs/flock"
	"github.com/kiyutink/tipam/core"
	"gopkg.in/yaml.v3"
)

const (
	DEFAULT_YAML_STATE_API_VERSION = 1
	PersistorLocalYaml             = "localyaml"
)

type YAMLState struct {
	APIVersion   int                 `yaml:"apiVersion"`
	Reservations map[string][]string `yaml:"reservations,omitempty"`
	Claims       map[string][]string `yaml:"claims,omitempty"`
}

func newEmptyYAMLState() *YAMLState {
	return &YAMLState{
		APIVersion:   DEFAULT_YAML_STATE_API_VERSION,
		Reservations: map[string][]string{},
		Claims:       map[string][]string{},
	}
}

type LocalYAMLPersistor struct {
	fileName string
	flock    *flock.Flock
}

func NewLocalYAMLPersistor(fileName string) *LocalYAMLPersistor {
	return &LocalYAMLPersistor{fileName: fileName, flock: flock.New(fileName)}
}

func (lyp *LocalYAMLPersistor) Persist(s *core.State) error {
	file, err := os.OpenFile(lyp.fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	yamlState := newEmptyYAMLState()

	for cidr, r := range s.Reservations {
		yamlState.Reservations[cidr] = r.Tags
	}

	// for cidr, c := range s.Claims {
	// }

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)

	err = encoder.Encode(yamlState)

	return err
}

func (lyp *LocalYAMLPersistor) Read() (*core.State, error) {
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

	state := &core.State{
		Reservations: map[string]core.Reservation{},
	}

	for c, tags := range yamlState.Reservations {
		_, ipNet, err := net.ParseCIDR(c)
		if err != nil {
			return nil, err
		}
		r := core.NewReservation(ipNet, tags)
		state.Reservations[c] = r
	}

	return state, nil
}

func (lyp *LocalYAMLPersistor) Lock() error {
	return lyp.flock.Lock()
}

func (lyp *LocalYAMLPersistor) Unlock() error {
	return lyp.flock.Unlock()
}
