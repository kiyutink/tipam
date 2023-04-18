package tipam

import (
	"fmt"
	"net"
)

func (r *Runner) Get(cidr string) error {
	_, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}

	err = r.persistor.Lock()
	if err != nil {
		return err
	}

	state, err := r.persistor.Read()
	if err != nil {
		return err
	}

	if res, ok := state.Claims[cidr]; ok {
		fmt.Printf("%v", res)
	}

	return nil
}
