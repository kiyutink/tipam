package tipam

type testPersistor struct {
	didLock   bool
	didUnlock bool
	state     *State
}

func (tp *testPersistor) Persist(s *State) error {
	tp.state = s
	return nil
}

func (tp *testPersistor) Read() (*State, error) {
	return tp.state, nil
}

func (tp *testPersistor) Lock() error {
	tp.didLock = true
	return nil
}

func (tp *testPersistor) Unlock() error {
	tp.didUnlock = true
	return nil
}

func newTestPersistor(state *State) *testPersistor {
	tp := &testPersistor{
		state: state,
	}
	return tp
}
