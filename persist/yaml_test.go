package persist

import (
	"testing"

	"github.com/kiyutink/tipam/tipam"
)

// Not sure if this test is needed, looks like it's testing the yaml package tbh
func TestStateToYAMLString(t *testing.T) {
	state := tipam.NewStateWithClaims([]*tipam.Claim{tipam.MustParseClaimFromCIDR("10.0.0.0/12", []string{"production"}, false)})
	expected := `apiVersion: 1
claims:
  10.0.0.0/12:
    tags:
      - production
`
	res, err := StateToYAMLString(state)
	if err != nil {
		t.Fatalf("got error calling StateToYAMLString: %v", err)
	}

	if res != expected {
		t.Fatalf("TestStateToYAMLString output doesn't match: expected %v, got %v", expected, res)
	}
}
