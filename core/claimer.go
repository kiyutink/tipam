package core

import (
	"fmt"
	"net"
)

type ClaimFlags struct {
	// If ComplySubs is set to true, subclaims of the newly created claim
	// will be made comply with the newly-created claim by prepending their
	// taglists with the tags of the new claim
	ComplySubs bool
}

func (r *Runner) Claim(cidr string, tags []string, flags ClaimFlags) error {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("error parsing cidr %v: %w", cidr, err)
	}
	newClaim := NewClaim(ipNet, tags)

	if r.doLock {
		r.persistor.Lock()
		defer r.persistor.Unlock()
	}

	state, err := r.persistor.Read()
	if err != nil {
		return fmt.Errorf("error reading state: %w", err)
	}

	subs, supers := state.FindRelated(newClaim)

	err = ValidateOnSupers(newClaim, supers)
	if flags.ComplySubs {
		for _, sub := range subs {
			tags = append([]string{}, newClaim.Tags...)
			tags = append(tags, sub.Tags...)

			sub.Tags = tags
			state.Claims[sub.IPNet.String()] = sub
		}
	} else {
		err = ValidateOnSubs(newClaim, subs)
	}

	if err != nil {
		return fmt.Errorf("the claim is invalid: %w", err)
	}

	state.Claims[newClaim.IPNet.String()] = newClaim
	err = r.persistor.Persist(state)
	if err != nil {
		return fmt.Errorf("error writing state: %w", err)
	}

	return nil
}
