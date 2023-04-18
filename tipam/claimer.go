package tipam

import (
	"fmt"
	"net"
)

type claimParams struct {
	complySubs bool
}
type ClaimOption func(*claimParams)

func WithComplySubs(v bool) ClaimOption {
	return func(opts *claimParams) {
		opts.complySubs = true
	}
}

func (r *Runner) Claim(cidr string, tags []string, o ...ClaimOption) error {
	params := &claimParams{}
	for _, opt := range o {
		opt(params)
	}
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
	if params.complySubs {
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
