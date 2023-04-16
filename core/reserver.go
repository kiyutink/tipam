package core

import (
	"fmt"
	"net"
)

func (r *Runner) Reserve(cidr string, tags []string) error {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("error parsing cidr %v: %w", cidr, err)
	}
	newReservation := NewReservation(ipNet, tags)

	r.Persistor.Lock()
	defer r.Persistor.Unlock()

	state, err := r.Persistor.Read()
	if err != nil {
		return fmt.Errorf("error reading state: %w", err)
	}

	err = state.ValidateReservation(newReservation)
	if err != nil {
		return fmt.Errorf("the reservation is invalid: %w", err)
	}

	state.Reservations[cidr] = newReservation
	err = r.Persistor.Persist(state)
	if err != nil {
		return fmt.Errorf("error writing state: %w", err)
	}

	return nil
}
