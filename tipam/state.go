package tipam

type State struct {
	Claims map[string]Claim
}

// FindRelated returns all the related Claims as subs, supers
func (s *State) FindRelated(cl Claim) ([]Claim, []Claim) {
	subs, supers := []Claim{}, []Claim{}
	for _, r := range s.Claims {
		if r.LiesWithinRangeOf(cl) {
			subs = append(subs, r)
		}
		if cl.LiesWithinRangeOf(r) {
			supers = append(supers, r)
		}
	}

	return subs, supers
}

func (s *State) FindSubs(cl Claim) []Claim {
	sub, _ := s.FindRelated(cl)
	return sub
}

func (s *State) FindSupers(cl Claim) []Claim {
	_, super := s.FindRelated(cl)
	return super
}
