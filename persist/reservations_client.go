package persist

import (
	"fmt"
	"net"
	"os"

	"github.com/kiyutink/tipam/core"
	"gopkg.in/yaml.v3"
)

type YamlReservationsClient struct{}

type reservationYamlRecord struct {
	CIDR string   `yaml:"cidr"`
	Tags []string `yaml:"tags"`
}

func (yrc *YamlReservationsClient) Create(reservation core.Reservation) error {
	file, err := os.OpenFile("testdata/reservations.yaml", os.O_WRONLY|os.O_APPEND, 0) // TODO: don't hardcode the filename
	defer file.Close()
	if err != nil {
		return fmt.Errorf("error persisting reservation to yaml: %w", err)
	}

	err = yaml.NewEncoder(file).Encode([]reservationYamlRecord{{
		CIDR: reservation.IPNet.String(),
		Tags: reservation.Tags,
	}})

	if err != nil {
		return fmt.Errorf("error persisting reservation to yaml: %w", err)
	}

	return nil
}

func (yrc *YamlReservationsClient) ReadAll() ([]core.Reservation, error) {
	reservationRecords := []reservationYamlRecord{}

	file, err := os.Open("testdata/reservations.yaml") // TODO: don't hardcode the filename
	if err != nil {
		return nil, fmt.Errorf("error opening persistence yaml file: %w", err)
	}

	err = yaml.NewDecoder(file).Decode(&reservationRecords)
	if err != nil {
		return nil, fmt.Errorf("error reading reservations from yaml: %w", err)
	}

	reservations := make([]core.Reservation, len(reservationRecords))

	for i, resRec := range reservationRecords {
		_, ipNet, err := net.ParseCIDR(resRec.CIDR)
		if err != nil {
			return nil, fmt.Errorf("error parsing cidr %v from a reservation: %w", resRec.CIDR, err)
		}
		reservations[i] = core.Reservation{
			IPNet: ipNet,
			Tags:  resRec.Tags,
		}
	}

	return reservations, nil
}
