package tipam

import "fmt"

// ValidateOnSubs validates (duh) claim r against given list of
// subclaims. All subclaims should have all the tags
// present on r and should have longer taglists
func ValidateOnSubs(c *Claim, subs []*Claim) error {
	for _, cl := range subs {
		if !cl.IsValidSubclaimOf(c) {
			return fmt.Errorf("the claim is not a valid superclaim of claim with CIDR=%v", cl.IPNet.String())
		}
	}
	return nil
}

// ValidateOnSupers validates (duh) claim r against a list of superclaims
// The claim r should have all the tags that the longest super has
// and introduce at least one new tag
func ValidateOnSupers(c *Claim, supers []*Claim) error {
	for _, cl := range supers {
		if !c.IsValidSubclaimOf(cl) {
			return fmt.Errorf("the claim is not a valid subclaim of claim with CIDR=%v", cl.IPNet.String())
		}
	}
	return nil
}
