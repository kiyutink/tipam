package persist

import (
	"bytes"
	"io"

	"github.com/kiyutink/tipam/tipam"
	"gopkg.in/yaml.v3"
)

const (
	defaultYAMLStateAPIVersion = 1
)

// yamlStateClaim is a yaml of tipam.Claim. Used for (un)marshaling
type yamlStateClaim struct {
	Tags  []string `yaml:"tags"`
	Final bool     `yaml:"final,omitempty"`
}

// yamlState is a YAML representation of tipam.State. Used for (un)marshaling
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

// StateToYAMLString converts a *tipam.State to a YAML string
func StateToYAMLString(s *tipam.State) (string, error) {
	buf := bytes.NewBuffer([]byte{})
	ys := stateToYAMLState(s)

	err := encodeYAMLState(ys, buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
