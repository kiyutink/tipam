package tipam

func (r *Runner) Release(cidr string) error {
	err := r.persistor.Lock()
	if err != nil {
		return err
	}
	defer r.persistor.Unlock()
	state, err := r.persistor.Read()
	if err != nil {
		return err
	}

	delete(state.Claims, cidr)
	err = r.persistor.Persist(state)

	return err
}
