package core

import "fmt"

// ValidateOnSubs validates (duh) reservation r against given list of
// subreservations. All subreservations should have all the tags
// present on r and should have longer taglists
func ValidateOnSubs(r Reservation, subs []Reservation) error {
	for _, res := range subs {
		if !res.IsValidSubreservationOf(r) {
			return fmt.Errorf("the reservation is not a valid subreservation of reservation with CIDR=%v", res.IPNet.String())
		}
	}
	return nil
}

// ValidateOnSupers validates (duh) reservation r against a list of superresrevations.
// The reservation r should have all the tags that the longest super has
// and introduce at least one new tag
func ValidateOnSupers(r Reservation, supers []Reservation) error {
	for _, res := range supers {
		if !r.IsValidSubreservationOf(res) {
			return fmt.Errorf("the reservation is not a valid superreservation of reservation with CIDR=%v", res.IPNet.String())
		}
	}
	return nil
}
