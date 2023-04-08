package persist

import (
	"fmt"
	"net"
	"os"

	"github.com/kiyutink/tipam/core"
	"gopkg.in/yaml.v3"
)

type LocalYamlReservationsClient struct {
	fileName string
}

func NewLocalYamlReservationsClient(fileName string) *LocalYamlReservationsClient {
	return &LocalYamlReservationsClient{
		fileName: fileName,
	}
}

type stateFile struct {
	APIVersion   int                 `yaml:"apiVersion"`
	Reservations map[string][]string `yaml:"reservations,omitempty"`
	Claims       map[string][]string `yaml:"claims,omitempty"`
}

func (yrc *LocalYamlReservationsClient) readState() (*stateFile, error) {
	bytes, err := os.ReadFile(yrc.fileName)
	if err != nil {
		return nil, fmt.Errorf("error persisting reservation to yaml: %w", err)
	}
	state := stateFile{}

	err = yaml.Unmarshal(bytes, &state)
	return &state, nil
}

func (yrc *LocalYamlReservationsClient) Create(reservation core.Reservation) error {
	state, err := yrc.readState()
	if err != nil {
		return fmt.Errorf("error persisting reservation to yaml: %w", err)
	}

	state.Reservations[reservation.IPNet.String()] = reservation.Tags

	file, err := os.OpenFile(yrc.fileName, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return fmt.Errorf("error persisting reservation to yaml: %w", err)
	}

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	err = encoder.Encode(state)
	if err != nil {
		return fmt.Errorf("error persisting reservation to yaml: %w", err)
	}

	return nil
}

func (yrc *LocalYamlReservationsClient) ReadAll() ([]core.Reservation, error) {
	state, err := yrc.readState()
	if err != nil {
		return nil, fmt.Errorf("error persisting reservation to yaml: %w", err)
	}

	reservations := []core.Reservation{}

	for CIDR, tags := range state.Reservations {
		_, ipNet, err := net.ParseCIDR(CIDR)
		if err != nil {
			return nil, fmt.Errorf("error parsing cidr %v from a reservation: %w", CIDR, err)
		}
		reservations = append(reservations, core.Reservation{
			IPNet: ipNet,
			Tags:  tags,
		})
	}

	return reservations, nil
}
