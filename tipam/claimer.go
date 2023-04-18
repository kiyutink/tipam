package tipam

import (
	"errors"
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

func (r *Runner) Claim(cidr string, tags []string, opts ...ClaimOption) error {
	params := &claimParams{
		complySubs: false,
	}
	for _, opt := range opts {
		opt(params)
	}

	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("error parsing cidr %v: %w", cidr, err)
	}

	if len(tags) < 1 {
		return fmt.Errorf("at least one tag has to be provided")
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

	if _, ok := state.Claims[newClaim.IPNet.String()]; ok {
		return errors.New("a claim for this CIDR already exists")
	}

	subs, supers := state.FindRelated(newClaim)

	err = ValidateOnSupers(newClaim, supers)

	if err != nil {
		return err
	}

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
