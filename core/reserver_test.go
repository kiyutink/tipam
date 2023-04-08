package core

import (
	"fmt"
	"testing"
)

func TestReserve(t *testing.T) {
	tr := Runner{
		Persistor: &testReservationsClient{
			testCreate: func(r Reservation) error {
				return nil
			},
			testReadAll: func() ([]Reservation, error) {
				r, _ := NewReservation("10.0.1.0/24", []string{"test", "test_inner"})
				return []Reservation{*r}, nil
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
