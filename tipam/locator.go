package tipam

import (
	"bytes"
	"errors"
	"net"
	"slices"
	"sort"

	gocidr "github.com/apparentlymart/go-cidr/cidr"
)

// Locate finds a spot for a claim and returns the cidr
func (r *Runner) Locate(within []string, maskLen int) (string, error) {
	err := r.persistor.Lock()
	if err != nil {
		return "", err
	}
	defer r.persistor.Unlock()

	state, err := r.persistor.Read()
	if err != nil {
		return "", err
	}

	potentialSupers := []*Claim{}

	// Out of all the claims in the state, find the ones that match by tags.
	for _, claim := range state.Claims {
		if !slices.Equal(claim.Tags, within) {
			continue
		}

		// We don't want to consider final claims as potential supers.
		if claim.Final {
			continue
		}

		claimNetMaskOnes, _ := claim.IPNet.Mask.Size()
		// This means the potential superclaim isn't large enough and can be skipped.
		if maskLen <= claimNetMaskOnes {
			continue
		}

		potentialSupers = append(potentialSupers, claim)
	}

	// We want to sort all the potential supers, as the order of the records in the state is random,
	// and we want to place the new claim as "top left" as possible in the address space.
	sort.Slice(potentialSupers, func(i, j int) bool {
		return bytes.Compare(potentialSupers[i].IPNet.IP, potentialSupers[j].IPNet.IP) < 0
	})

	for _, ps := range potentialSupers {
		ipNet := &net.IPNet{
			IP: ps.IPNet.IP,
			// We assume here that we're only working with 32-bit masks!
			Mask: net.CIDRMask(maskLen, 32),
		}

		for ps.IPNet.Contains(ipNet.IP) {
			// We will create a dummy dummyClaim to check if it is valid against the state.
			dummyClaim := newDummyClaim(ipNet, within)
			subs, supers := state.FindRelated(dummyClaim)
			validAgainstSupers := dummyClaim.ValidateSupers(supers) == nil
			validAgainstSubs := dummyClaim.ValidateSubs(subs) == nil
			if _, ok := state.Claims[dummyClaim.IPNet.String()]; !ok && validAgainstSupers && validAgainstSubs {
				// We found a place for the claim, return it.
				return ipNet.String(), nil
			}
			ipNet, _ = gocidr.NextSubnet(ipNet, maskLen)
		}
	}

	return "", errors.New("no possible network found")
}

func newDummyClaim(ipNet *net.IPNet, tags []string) *Claim {
	dummyTags := make([]string, len(tags))
	copy(dummyTags, tags)
	dummyTags = append(dummyTags, "dummy")
	return NewClaim(ipNet, dummyTags, false)
}
