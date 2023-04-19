package cmd

import (
	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	var cidr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "show claim with given CIDR",
		Long:  "show claim with given CIDR, if it exists. If it doesn't exist, do nothing",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := newPersistor()
			if err != nil {
				return err
			}
			runner := newRunner(p)

			return runner.Get(cidr)
		},
	}

	cmd.Flags().StringVar(&cidr, "cidr", "", "the cidr of the claim to be removed")

	return cmd
}
