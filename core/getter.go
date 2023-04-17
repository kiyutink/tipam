package core

import "fmt"

func (r *Runner) Get(cidr string) error {
	err := r.persistor.Lock()
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
