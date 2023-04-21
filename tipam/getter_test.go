package tipam

import "testing"

func TestGetter(t *testing.T) {
	testCases := []struct {
		cidr     string
		succeeds bool
	}{
		{"10.0.0.0/8", true},
		{"172.0.0.0/8", false},
	}

	for _, tc := range testCases {
		claims := []*Claim{
			MustParseClaimFromCIDR("10.0.0.0/8", []string{"test"}, false),
			MustParseClaimFromCIDR("20.0.0.0/8", []string{"test"}, false),
		}

		runner := NewRunner(newTestPersistor(NewStateWithClaims(claims)), RunnerOpts{})
		cl, err := runner.Get(tc.cidr)

		if tc.succeeds {
			if err != nil {
				t.Errorf("expected get to succeed but received a non-nil error")
			}
			if cl == nil {
				t.Errorf("expected get to return a claim, but received nil")
			}
			if cl.IPNet.String() != tc.cidr {
				t.Errorf("claim returned by get doesn't match the claim in the state")
			}
		} else {
			if err == nil {
				t.Errorf("expected get to fail but received a nil error")
			}

			if cl != nil {
				t.Errorf("expected claim to be nil, but received %T", cl)
			}
		}
	}
}
