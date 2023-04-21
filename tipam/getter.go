package tipam

import (
	"errors"
	"net"
)

func (r *Runner) Get(cidr string) (*Claim, error) {
	_, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	err = r.persistor.Lock()
	if err != nil {
		return nil, err
	}
	defer r.persistor.Unlock()

	state, err := r.persistor.Read()
	if err != nil {
		return nil, err
	}

	if cl, ok := state.Claims[cidr]; ok {
		return cl, nil
	}
	return nil, errors.New("claim not found")
}
