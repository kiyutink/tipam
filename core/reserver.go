package core

import (
	"fmt"
	"net"
)

type ReserveFlags struct {
	// If ComplySubs is set to true, subreservations of the newly created reservation
	// will be made comply with the newly-created reservation by prepending their
	// taglists with the tags of the new reservation
	ComplySubs bool
}

func (r *Runner) Reserve(cidr string, tags []string, flags ReserveFlags) error {
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

	subs, supers := state.FindRelated(newReservation)

	err = ValidateOnSupers(newReservation, supers)
	if flags.ComplySubs {
		for _, sub := range subs {
			tags = append([]string{}, newReservation.Tags...)
			tags = append(tags, sub.Tags...)

			sub.Tags = tags
			state.Reservations[sub.IPNet.String()] = sub
		}
	} else {
		err = ValidateOnSubs(newReservation, subs)
	}

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
