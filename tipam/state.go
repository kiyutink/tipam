package tipam

import "fmt"

type State struct {
	Claims map[string]Claim
}

func (s *State) ValidateClaim(newRes Claim) error {
	for _, res := range s.Claims {
		if !newRes.LiesWithinRangeOf(res) {
			continue
		}

		if !newRes.IsValidSubclaimOf(res) {
			return fmt.Errorf("the claim is not a valid subclaim of claim with CIDR=%v", res.IPNet.String())
		}
	}
	return nil
}

func (s *State) FindRelated(res Claim) ([]Claim, []Claim) {
	subs, supers := []Claim{}, []Claim{}
	for _, r := range s.Claims {
		if r.LiesWithinRangeOf(res) {
			subs = append(subs, r)
		}
		if res.LiesWithinRangeOf(r) {
			supers = append(supers, r)
		}
	}

	return subs, supers
}

func (s *State) FindSubs(res Claim) []Claim {
	sub, _ := s.FindRelated(res)
	return sub
}

func (s *State) FindSupers(res Claim) []Claim {
	_, super := s.FindRelated(res)
	return super
}
