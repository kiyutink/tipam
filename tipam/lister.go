package tipam

func (r *Runner) List() (*State, error) {
	state, err := r.persistor.Read()
	if err != nil {
		return nil, err
	}

	return state, nil
}
