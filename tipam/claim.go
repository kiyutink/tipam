package tipam

import "net"

type Claim struct {
	IPNet *net.IPNet
	Tags  []string
	Final bool
}

// NewClaim returns a new Claim given its fields
func NewClaim(ipNet *net.IPNet, tags []string, final bool) *Claim {
	return &Claim{
		IPNet: ipNet,
		Tags:  tags,
		Final: final,
	}
}

// ParseClaimFromCIDR parses a given CIDR and returns a Claim and an optional error
func ParseClaimFromCIDR(cidr string, tags []string, final bool) (*Claim, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	return NewClaim(ipNet, tags, final), nil
}

// MustParseClaimFromCIDR parses a given CIDR and returns a Claim. Panics if the CIDR is invalid
func MustParseClaimFromCIDR(cidr string, tags []string, final bool) *Claim {
	c, _ := ParseClaimFromCIDR(cidr, tags, final)
	return c
}

// LiesWithinRangeOf checks whether the claim's CIDR range lies within the range of the supernet
func (c *Claim) LiesWithinRangeOf(super *Claim) bool {
	onesSub, _ := c.IPNet.Mask.Size()
	onesSuper, _ := super.IPNet.Mask.Size()

	// The subnet's mask should be longer, constituting a smaller CIDR range
	if onesSub <= onesSuper {
		return false
	}

	// The subnet's network address must lie within the supernet's CIDR range
	if !super.IPNet.Contains(c.IPNet.IP) {
		return false
	}

	return true
}

// IsValidSubclaimOf checks whether the claim has all the necessary tags to be a valid subclaim of super
func (r *Claim) IsValidSubclaimOf(super *Claim) bool {
	// The subclaim must introduce at least one new tag.
	if len(r.Tags) <= len(super.Tags) {
		return false
	}

	// You can't create a subclaim of a final claim
	if super.Final {
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
