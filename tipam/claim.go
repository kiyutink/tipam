package tipam

import "net"

type Claim struct {
	IPNet *net.IPNet
	Tags  []string
	Final bool
}

// NewClaimFromCIDR parses a given CIDR and returns a claim. If the CIDR is invalid, returns an error
func NewClaim(ipNet *net.IPNet, tags []string) Claim {
	return Claim{
		IPNet: ipNet,
		Tags:  tags,
	}
}

// LiesWithinRangeOf checks whether the claim's CIDR range lies within the range of the supernet
func (c Claim) LiesWithinRangeOf(super Claim) bool {
	onesSub, _ := c.IPNet.Mask.Size()
	onesSuper, _ := super.IPNet.Mask.Size()

	// The subnet's mask should be longer, constituting a smaller CIDR range
	if onesSub < onesSuper {
		return false
	}

	// The subnet's network address must lie within the supernet's CIDR range
	if !super.IPNet.Contains(c.IPNet.IP) {
		return false
	}

	return true
}

// IsValidSubclaimOf checks whether the claim has all the necessary tags to be a valid subclaim of super
func (r Claim) IsValidSubclaimOf(super Claim) bool {
	// The subclaim must introduce at least one new tag.
	if len(r.Tags) <= len(super.Tags) {
		return false
	}

	// The subclaim must have all the tags that the superclaim has. The tags have to be in the same order
	for i := range super.Tags {
		if super.Tags[i] != r.Tags[i] {
			return false
		}
	}

	return true
}
