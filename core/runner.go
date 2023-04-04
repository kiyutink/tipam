package core

type reservationsClient interface {
	Create(reservation Reservation) error
	ReadAll() ([]Reservation, error)
	// Replace(cidr string, reservation Reservation) error
	// Delete(cidr string) error
}

type Runner struct {
	ReservationsClient reservationsClient
}
