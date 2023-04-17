package core

import "fmt"

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

func (s *State) FindRelated(res Reservation) ([]Reservation, []Reservation) {
	subs, supers := []Reservation{}, []Reservation{}
	for _, r := range s.Reservations {
		if r.LiesWithinRangeOf(res) {
			subs = append(subs, r)
		}
		if res.LiesWithinRangeOf(r) {
			supers = append(supers, r)
		}
	}

	return subs, supers
}

func (s *State) FindSubs(res Reservation) []Reservation {
	sub, _ := s.FindRelated(res)
	return sub
}

func (s *State) FindSupers(res Reservation) []Reservation {
	_, super := s.FindRelated(res)
	return super
}
