package core

import "net"

type Reservation struct {
	IPNet *net.IPNet
	Tags  []string
}

// LiesWithinRangeOf checks whether the reservation's CIDR range lies within the range of parent
func (r *Reservation) LiesWithinRangeOf(parent Reservation) bool {
	onesChild, _ := r.IPNet.Mask.Size()
	onesParent, _ := parent.IPNet.Mask.Size()

	// The child block's mask should be longer, constituting a smaller CIDR range
	if onesChild < onesParent {
		return false
	}

	// The child resrevation's network address must lie within the parent's CIDR range
	if !parent.IPNet.Contains(r.IPNet.IP) {
		return false
	}

	return true
}

// IsValidSubreservationOf checks whether the reservation has all the necessary tags to be a valid subreservation of parent
func (r *Reservation) IsValidSubreservationOf(parent Reservation) bool {
	// The subreservation must introduce at least one new tag.
	if len(r.Tags) <= len(parent.Tags) {
		return false
	}

	// The subreservation must have all the tags that the "parent" reservation has. The tags have to be in the same order
	for i := range parent.Tags {
		if parent.Tags[i] != r.Tags[i] {
			return false
		}
	}

	return true
}
