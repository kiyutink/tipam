package core

func (r *Runner) Release(cidr string) error {
	state, err := r.Persistor.Read()
	if err != nil {
		return err
	}

	delete(state.Reservations, cidr)
	err = r.Persistor.Persist(state)

	return err
}
