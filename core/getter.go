package core

// TODO: This seems a little weirdly specific and probably deserves some more thinking.
// A good start though, we need some way of loading this into a map
func (r *Runner) GetTags() (map[string][]string, error) {
	reservations, err := r.ReservationsClient.ReadAll()
	if err != nil {
		return nil, err
	}

	m := map[string][]string{}
	for _, res := range reservations {
		m[res.IPNet.String()] = res.Tags
	}

	return m, nil
}
