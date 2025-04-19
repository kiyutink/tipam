package tipam

import (
	"fmt"
	"net"
)

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
	return isSubnet(c.IPNet, super.IPNet)
}

// IsValidSubclaimOf checks whether the claim has all the necessary tags to be a valid subclaim of super
func (c *Claim) IsValidSubclaimOf(super *Claim) bool {
	// The subclaim must introduce at least one new tag.
	if len(c.Tags) <= len(super.Tags) {
		return false
	}

	// You can't create a subclaim of a final claim
	if super.Final {
		return false
	}

	// The subclaim must have all the tags that the superclaim has. The tags have to be in the same order
	for i := range super.Tags {
		if super.Tags[i] != c.Tags[i] {
			return false
		}
	}

	return true
}

// ValidateSubs validates (duh) claim r against given list of subclaims
// subclaims. All subclaims should have all the tags
// present on r and should have longer taglists
func (c *Claim) ValidateSubs(subs []*Claim) error {
	for _, cl := range subs {
		if !cl.IsValidSubclaimOf(c) {
			return fmt.Errorf("the claim is not a valid superclaim of claim with CIDR=%v", cl.IPNet.String())
		}
	}
	return nil
}

// ValidateSupers validates (duh) claim r against a list of superclaims
// The claim r should have all the tags that the longest super has
// and introduce at least one new tag
func (c *Claim) ValidateSupers(supers []*Claim) error {
	for _, cl := range supers {
		if !c.IsValidSubclaimOf(cl) {
			return fmt.Errorf("the claim is not a valid subclaim of claim with CIDR=%v", cl.IPNet.String())
		}
	}
	return nil
}
