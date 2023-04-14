package core

import (
	"fmt"
)

type State struct {
	Reservations map[string]Reservation
	Claims       map[string]Claim
}

func (s *State) ValidateReservation(newRes Reservation) error {
	for _, res := range s.Reservations {
		if !newRes.LiesWithinRangeOf(res) {
			continue
		}

		if !newRes.IsValidSubreservationOf(res) {
			return fmt.Errorf("the reservation is not a valid subreservation of reservation with CIDR=%v", res.IPNet.String())
		}
	}
	return nil
}

func (s *State) FindParentReservations(res Reservation) []Reservation {
	reservations := []Reservation{}
	for _, r := range s.Reservations {
		if res.LiesWithinRangeOf(r) {
			reservations = append(reservations, r)
		}
	}

	return reservations
}
