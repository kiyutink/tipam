package tipam

import (
	"fmt"
	"net"
	"testing"
)

func TestClaim(t *testing.T) {
	tr := Runner{
		persistor: &testPersistor{
			testPersist: func(s *State) error {
				return nil
			},
			testRead: func() (*State, error) {
				cidr := "10.0.1.0/24"
				_, ipNet, _ := net.ParseCIDR(cidr)
				r := NewClaim(ipNet, []string{"test", "test_inner"})
				return &State{
					Claims: map[string]Claim{
						cidr: r,
					},
				}, nil
			},
		},
	}

	tests := []struct {
		cidr  string
		tags  []string
		fails bool
	}{
		{"10.0.0.0/24", []string{"test2"}, false},
		{"10.0.1.0/25", []string{"test2"}, true},
		{"10.0.1.0/25", []string{"test"}, true},
		{"10.0.1.0/25", []string{"test", "test_inner"}, true},
		{"10.0.1.0/25", []string{"test", "test_inner", "test_inner_2"}, false},
	}

	for _, tt := range tests {
		err := tr.Claim(tt.cidr, tt.tags)
		if tt.fails && err == nil {
			t.Errorf("expected claim to fail, but didn't")
		}

		if !tt.fails && err != nil {
			t.Errorf("expected claim to succeed, but didn't")
		}
	}

	fmt.Println(tr)
}
