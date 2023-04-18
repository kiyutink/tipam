package tipam

type testPersistor struct {
	Persistor
	testPersist func(*State) error
	testRead    func() (*State, error)
}

func (tp *testPersistor) Persist(s *State) error {
	return tp.testPersist(s)
}

func (tp *testPersistor) Read() (*State, error) {
	return tp.testRead()
}
