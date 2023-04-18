package tipam

import (
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
		{"10.0.1.0/25", []string{"test2"}, true, "should fail if not enough tags"},
		{"10.0.1.0/25", []string{"test"}, true, "should fail if not enough tags"},
		{"10.0.1.0/25", []string{"test", "test_inner"}, true, "should fail if not enough tags"},
		{"10.0.1.0/20", []string{"test_outer"}, true, "should fail if trying to claim a larger range without complySubs"},

		{"10.0.0.0/24", []string{"test2"}, false, "should succeed if the CIDR is free"},
		{"10.0.1.0/25", []string{"test", "test_inner", "test_inner_2"}, false, "should succeed if valid"},
	}

	for _, tc := range testCases {
		tr := Runner{
			persistor: newTestPersistor(),
		}

		err := tr.Claim(tc.cidr, tc.tags)

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
		didLock := false
		didUnlock := false

		tp := newTestPersistor()
		tp.testLock = func() error { didLock = true; return nil }
		tp.testUnlock = func() error { didUnlock = true; return nil }

		tr := NewRunner(tp, WithLocking(tc.doLock))

		tr.Claim("172.16.0.0/12", []string{"test"})

		if didLock && !didUnlock {
			t.Errorf("state locked but not unlocked")
		}

		if !didLock && didUnlock {
			t.Errorf("state not locked, but tried to unlock")
		}

		if tc.doLock != didLock {
			t.Error(tc.errorMsg)
		}
	}
}

func TestClaimComplySubs(t *testing.T) {
	// TODO:
}
