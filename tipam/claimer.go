package tipam

import (
	"errors"
	"fmt"
)

// ClaimOpts is a struct of options (or flags) for the Claim method
type ClaimOpts struct {
	// ComplySubs reflects whether the subclaims of a newly-created claim
	// will be made comply (read: prepended with the list of tags) with the new claim.
	// r.Claim will error out if the claim has subclaims, and ComplySubs is set to `false`.
	// The default is `false`.
	ComplySubs bool
}

// Claim creates a new claim
func (r *Runner) Claim(c *Claim, opts ClaimOpts) error {
	if len(c.Tags) < 1 {
		return fmt.Errorf("at least one tag has to be provided")
	}

	for _, tag := range c.Tags {
		if tag == "" {
			return fmt.Errorf("tag can not be an empty string")
		}
	}

	if r.doLock {
		r.persistor.Lock()
		defer r.persistor.Unlock()
	}

	state, err := r.persistor.Read()
	if err != nil {
		return fmt.Errorf("error reading state: %w", err)
	}

	if _, ok := state.Claims[c.IPNet.String()]; ok {
		return errors.New("a claim for this CIDR already exists")
	}

	subs, supers := state.FindRelated(c)

	err = ValidateOnSupers(c, supers)

	if err != nil {
		return err
	}

	if opts.ComplySubs {
		if c.Final {
			err = errors.New("can't create a final superclaim")
		}
		for _, sub := range subs {
			sub.Tags = append(c.Tags[0:len(c.Tags):len(c.Tags)], sub.Tags...)
		}
	} else {
		err = ValidateOnSubs(c, subs)
	}

	if err != nil {
		return fmt.Errorf("the claim is invalid: %w", err)
	}

	state.Claims[c.IPNet.String()] = c
	err = r.persistor.Persist(state)
	if err != nil {
		return fmt.Errorf("error writing state: %w", err)
	}

	return nil
}
