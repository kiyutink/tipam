package core

type Persistor interface {
	Create(reservation Reservation) error
	ReadAll() ([]Reservation, error)
	// Replace(cidr string, reservation Reservation) error
	// Delete(cidr string) error
}

type Runner struct {
	Persistor Persistor
}

func NewRunner(p Persistor) *Runner {
	return &Runner{
		Persistor: p,
	}
}
