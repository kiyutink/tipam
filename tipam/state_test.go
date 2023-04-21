package tipam

import (
	"testing"
)

func TestFindRelated(t *testing.T) {
	cidrs := []string{
		"10.0.0.0/12",
		"10.0.0.0/16",
		"10.1.0.0/16",
		"10.0.0.7/32",
		"172.0.0.0/30",
	}

	claims := []*Claim{}
	for _, cidr := range cidrs {
		claims = append(claims, MustParseClaimFromCIDR(cidr, []string{"test"}, false))
	}
	state := NewStateWithClaims(claims)

	subs, supers := state.FindRelated(MustParseClaimFromCIDR("10.0.0.0/20", []string{"test"}, false))

	if len(subs) != 1 {
		t.Errorf("expected 1 subclaims, instead got %v", len(subs))
	}

	if len(supers) != 2 {
		t.Errorf("expected 2 superclaims, instead got %v", len(supers))
	}
}
