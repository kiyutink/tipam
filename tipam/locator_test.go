package tipam

import (
	"testing"
)

func TestLocate(t *testing.T) {
	// Should be able to locate when there's an appropriate place
	testCases := []struct {
		tags         []string
		maskLen      int
		expectedCidr string
		fails        bool
		errorMsg     string
	}{
		{[]string{"test"}, 16, "10.6.0.0/16", false, "Should successfully locate"},
		{[]string{"test", "test_inner"}, 24, "10.0.0.0/24", false, "Should successfully locate with 2 tags"},
		{[]string{"test", "test_inner_2"}, 24, "10.2.0.0/24", false, "Should successfully locate left-most"},
		{[]string{"test", "test_inner"}, 6, "", true, "Should not be able to locate when maskLen less than largest superclaim"},
		{[]string{"nonexistent"}, 24, "", true, "Should not be able to locate when no fitting superclaim exists"},
	}

	for _, tc := range testCases {
		existingClaims := []*Claim{
			MustParseClaimFromCIDR("10.0.0.0/8", []string{"test"}, false),
			MustParseClaimFromCIDR("10.0.0.0/16", []string{"test", "test_inner"}, false),
			MustParseClaimFromCIDR("10.1.0.0/16", []string{"test", "test_inner_2"}, true),
			MustParseClaimFromCIDR("10.2.0.0/16", []string{"test", "test_inner_2"}, false),
			MustParseClaimFromCIDR("10.3.0.0/16", []string{"test", "test_inner_2"}, false),
			MustParseClaimFromCIDR("10.4.0.0/16", []string{"test", "test_inner_3"}, true),
			MustParseClaimFromCIDR("10.5.0.0/16", []string{"test", "test_inner_3"}, false),
		}
		state := NewStateWithClaims(existingClaims)

		tr := Runner{
			persistor: newTestPersistor(state),
		}

		locatedCidr, err := tr.Locate(tc.tags, tc.maskLen)
		if err != nil && tc.fails {
			continue
		}

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if locatedCidr != tc.expectedCidr {
			t.Fatalf(tc.errorMsg)
		}
	}
}
