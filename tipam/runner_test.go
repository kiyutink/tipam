package tipam

import (
	"net"
)

type testPersistor struct {
	Persistor
	testPersist func(*State) error
	testRead    func() (*State, error)
	testLock    func() error
	testUnlock  func() error
}

func (tp *testPersistor) Persist(s *State) error {
	return tp.testPersist(s)
}

func (tp *testPersistor) Read() (*State, error) {
	return tp.testRead()
}

func (tp *testPersistor) Lock() error {
	return tp.testLock()
}

func (tp *testPersistor) Unlock() error {
	return tp.testUnlock()
}

func newTestPersistor() *testPersistor {
	tp := &testPersistor{
		testPersist: func(s *State) error {
			return nil
		},
		testRead: func() (*State, error) {
			cidr := "10.0.1.0/24"
			_, ipNet, _ := net.ParseCIDR(cidr)
			c := NewClaim(ipNet, []string{"test", "test_inner"}, false)

			return &State{
				Claims: map[string]Claim{
					cidr: c,
				},
			}, nil
		},
	}
	return tp
}
