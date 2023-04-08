package core

type testReservationsClient struct {
	Persistor
	testCreate  func(Reservation) error
	testReadAll func() ([]Reservation, error)
}

func (trc *testReservationsClient) Create(r Reservation) error {
	return trc.testCreate(r)
}

func (trc *testReservationsClient) ReadAll() ([]Reservation, error) {
	return trc.testReadAll()
}
