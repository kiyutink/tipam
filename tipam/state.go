package tipam

type State struct {
	Claims map[string]*Claim
}

func NewState() *State {
	return &State{
		Claims: map[string]*Claim{},
	}
}

func NewStateWithClaims(claims []*Claim) *State {
	state := NewState()
	for _, c := range claims {
		state.Claims[c.IPNet.String()] = c
	}

	return state
}

// FindRelated returns all the related Claims as subs, supers
func (s *State) FindRelated(cl *Claim) ([]*Claim, []*Claim) {
	subs, supers := []*Claim{}, []*Claim{}
	for _, c := range s.Claims {
		if c.LiesWithinRangeOf(cl) {
			subs = append(subs, c)
		}
		if cl.LiesWithinRangeOf(c) {
			supers = append(supers, c)
		}
	}

	return subs, supers
}

func (s *State) FindSubs(cl *Claim) []*Claim {
	sub, _ := s.FindRelated(cl)
	return sub
}

func (s *State) FindSupers(cl *Claim) []*Claim {
	_, super := s.FindRelated(cl)
	return super
}
