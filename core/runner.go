package core

import (
	"fmt"
	"net"
)

type reservationsClient interface {
	Create(reservation Reservation) error
	ReadAll() ([]Reservation, error)
	// Replace(cidr string, reservation Reservation) error
	// Delete(cidr string) error
}

type Runner struct {
	ReservationsClient reservationsClient
}

func (r *Runner) CreateReservation(cidr string, tags []string) error {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("error parsing cidr %v: %w", cidr, err)
	}
	reservation := Reservation{
		IPNet: ipNet,
		Tags:  tags,
	}
	allReservations, err := r.ReservationsClient.ReadAll()
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

	err = r.ReservationsClient.Create(reservation)
	if err != nil {
		return fmt.Errorf("error persisting reservation: %w", err)
	}

	return nil
}
