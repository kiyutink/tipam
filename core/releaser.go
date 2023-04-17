package core

func (r *Runner) Release(cidr string) error {
	state, err := r.persistor.Read()
	if err != nil {
		return err
	}

	delete(state.Claims, cidr)
	err = r.persistor.Persist(state)

	return err
}
