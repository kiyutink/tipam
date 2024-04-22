package tipam

import (
	"fmt"
	"testing"
)

func TestClaimValidates(t *testing.T) {
	testCases := []struct {
		cidr       string
		tags       []string
		shouldFail bool
		errorMsg   string
	}{
		{"10.0.1.0/25", []string{}, true, "should fail if no tags are provided"},
		{"10.0.1.0/25", []string{""}, true, "should fail when providing an empty string as tag"},
		{"10.0.1.0/25", []string{"", ""}, true, "should fail when providing an empty string as tag"},
		{"10.0.1.0/25", []string{"", "", "test"}, true, "should fail when providing an empty string as tag"},
		{"10.0.1.0/25", []string{"test2"}, true, "should fail if not enough tags"},
		{"10.0.1.0/25", []string{"test"}, true, "should fail if not enough tags"},
		{"10.0.1.0/25", []string{"test", "test_inner"}, true, "should fail if not enough tags"},
		{"10.0.1.0/20", []string{"test_outer"}, true, "should fail if trying to claim a larger range without complySubs"},

		{"10.0.0.0/24", []string{"test2"}, false, "should succeed if the CIDR is free"},
		{"10.0.1.0/25", []string{"test", "test_inner", "test_inner_2"}, false, "should succeed if valid"},
	}

	for _, tc := range testCases {
		existingClaim := MustParseClaimFromCIDR("10.0.1.0/24", []string{"test", "test_inner"}, false)
		state := NewStateWithClaims([]*Claim{existingClaim})

		tr := Runner{
			persistor: newTestPersistor(state),
		}

		err := tr.Claim(MustParseClaimFromCIDR(tc.cidr, tc.tags, false), ClaimOpts{})

		if tc.shouldFail && err == nil {
			t.Errorf("expected claim with CIDR %v and tags %v to fail, but didn't: %v", tc.cidr, tc.tags, tc.errorMsg)
		}

		if !tc.shouldFail && err != nil {
			t.Errorf("expected claim with CIDR %v and tags %v to succeed, but didn't: %v", tc.cidr, tc.tags, tc.errorMsg)
		}
	}
}

func TestClaimLocksState(t *testing.T) {
	testCases := []struct {
		doLock   bool
		errorMsg string
	}{
		{true, "expected to use locking, but didn't"},
		{false, "expected not to use locking, but did"},
	}
	for _, tc := range testCases {
		tp := newTestPersistor(NewState())
		tr := NewRunner(tp, RunnerOpts{DoLock: tc.doLock})

		tr.Claim(MustParseClaimFromCIDR("172.16.0.0/12", []string{"test"}, false), ClaimOpts{})

		if tp.didLock && !tp.didUnlock {
			t.Errorf("state locked but not unlocked")
		}

		if !tp.didLock && tp.didUnlock {
			t.Errorf("state not locked, but tried to unlock")
		}

		if tc.doLock != tp.didLock {
			t.Error(tc.errorMsg)
		}
	}
}

func TestClaimFailsSubclaimsOfFinal(t *testing.T) {
	existingClaim := MustParseClaimFromCIDR("10.0.1.0/24", []string{"test1"}, true)
	state := NewStateWithClaims([]*Claim{existingClaim})
	tr := Runner{persistor: newTestPersistor(state)}

	err := tr.Claim(MustParseClaimFromCIDR("10.0.1.0/30", []string{"test", "test_inner"}, false), ClaimOpts{})

	if err == nil {
		t.Errorf("expected to fail when creating a subclaim of a final claim, but succeeded")
	}
}

func TestClaimComplySubs(t *testing.T) {
	testCases := []struct {
		final    bool
		succeeds bool
	}{
		{true, false},
		{false, true},
	}

	for _, tc := range testCases {
		existingClaims := []*Claim{
			MustParseClaimFromCIDR("10.0.1.0/24", []string{"test1"}, true),
			MustParseClaimFromCIDR("10.0.2.0/24", []string{"test2"}, true),
		}
		state := NewStateWithClaims(existingClaims)
		tr := Runner{persistor: newTestPersistor(state)}

		err := tr.Claim(MustParseClaimFromCIDR("10.0.0.0/16", []string{"test_outer"}, tc.final), ClaimOpts{ComplySubs: true})

		if err == nil && !tc.succeeds {
			t.Errorf("expected to fail creating a superclaim marked as final, but succeeded")
		}

		if err != nil && tc.succeeds {
			t.Errorf("expected to succeed creating a superclaim marked as final, but failed")
		}

		if tc.succeeds && len(state.Claims["10.0.1.0/24"].Tags) != 2 {
			fmt.Printf("%+v\n", state.Claims["10.0.1.0/24"])
			fmt.Printf("%+v\n", state.Claims["10.0.2.0/24"])
			t.Errorf("expected claim with complySubs to prepend new tags to subclaims, but didn't")
		}
	}
}
