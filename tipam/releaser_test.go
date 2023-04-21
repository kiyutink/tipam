package tipam

import "testing"

func TestReleaser(t *testing.T) {
	cidr := "10.0.0.0/8"
	c := MustParseClaimFromCIDR(cidr, []string{"test"}, false)
	state := NewStateWithClaims([]*Claim{c})
	p := newTestPersistor(state)

	tr := NewRunner(p, RunnerOpts{})
	err := tr.Release(cidr)
	if err != nil {
		t.Errorf("error running release on %v: %v", cidr, err)
	}

	if _, ok := state.Claims[cidr]; ok {
		t.Errorf("expected release to remove claim from state, but didn't")
	}
}
