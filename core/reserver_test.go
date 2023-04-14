package core

import (
	"fmt"
	"net"
	"testing"
)

func TestReserve(t *testing.T) {
	tr := Runner{
		Persistor: &testPersistor{
			testPersist: func(s *State) error {
				return nil
			},
			testRead: func() (*State, error) {
				cidr := "10.0.1.0/24"
				_, ipNet, _ := net.ParseCIDR(cidr)
				r := NewReservation(ipNet, []string{"test", "test_inner"})
				return &State{
					Reservations: map[string]Reservation{
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
		err := tr.Reserve(tt.cidr, tt.tags)
		if tt.fails && err == nil {
			t.Errorf("expected reservation to fail, but didn't")
		}

		if !tt.fails && err != nil {
			t.Errorf("expected reservation to succeed, but didn't")
		}
	}

	fmt.Println(tr)
}
