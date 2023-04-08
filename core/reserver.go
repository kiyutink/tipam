package core

import (
	"fmt"
)

func (r *Runner) Reserve(cidr string, tags []string) error {
	reservation, err := NewReservation(cidr, tags)
	if err != nil {
		return fmt.Errorf("error parsing cidr %v: %w", cidr, err)
	}
	allReservations, err := r.Persistor.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading existing reservations: %w", err)
	}
	for _, existingReservation := range allReservations {
		if !reservation.LiesWithinRangeOf(existingReservation) {
			continue
		}

		if !reservation.IsValidSubreservationOf(existingReservation) {
			return fmt.Errorf("the reservation is not a valid subreservation of reservation with CIDR=%v", existingReservation.IPNet.String())
		}
	}

	err = r.Persistor.Create(*reservation)
	if err != nil {
		return fmt.Errorf("error persisting reservation: %w", err)
	}

	return nil
}
