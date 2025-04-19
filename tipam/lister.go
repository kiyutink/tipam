package tipam

func (r *Runner) List() (*State, error) {
	err := r.persistor.Lock()
	if err != nil {
		return nil, err
	}
	defer r.persistor.Unlock()
	state, err := r.persistor.Read()
	if err != nil {
		return nil, err
	}

	return state, nil
}
