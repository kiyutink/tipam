package tipam

import (
	"errors"
	"fmt"
	"net"
)

type ClaimOpts struct {
	ComplySubs bool
}

func (r *Runner) Claim(cidr string, tags []string, final bool, opts ClaimOpts) error {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("error parsing cidr %v: %w", cidr, err)
	}

	if len(tags) < 1 {
		return fmt.Errorf("at least one tag has to be provided")
	}

	newClaim := NewClaim(ipNet, tags, final)

	if r.doLock {
		r.persistor.Lock()
		defer r.persistor.Unlock()
	}

	state, err := r.persistor.Read()
	if err != nil {
		return fmt.Errorf("error reading state: %w", err)
	}

	if _, ok := state.Claims[newClaim.IPNet.String()]; ok {
		return errors.New("a claim for this CIDR already exists")
	}

	subs, supers := state.FindRelated(newClaim)

	err = ValidateOnSupers(newClaim, supers)

	if err != nil {
		return err
	}

	if opts.ComplySubs {
		if final {
			err = errors.New("can't create a final superclaim")
		}
		for _, sub := range subs {
			sub.Tags = append(newClaim.Tags[0:len(newClaim.Tags):len(newClaim.Tags)], sub.Tags...)
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
