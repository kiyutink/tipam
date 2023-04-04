package core

import "fmt"

func (r *Runner) Release(cidr string) error {
	fmt.Println("cidr to release:", cidr)
	return nil
}
