package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	var cidr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Show claim with given CIDR",
		Long:  "Show claim with given CIDR, if it exists. If it doesn't exist, do nothing",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := newPersistor()
			if err != nil {
				return err
			}
			runner := newRunner(p)

			claim, err := runner.Get(cidr)
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", claim)
			return nil
		},
	}

	cmd.Flags().StringVar(&cidr, "cidr", "", "the cidr of the claim to be removed")

	return cmd
}
