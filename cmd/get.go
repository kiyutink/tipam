package cmd

import (
	"github.com/kiyutink/tipam/tipam"
	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	var cidr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Show resrvation with given CIDR",
		Long:  "Show claim with given CIDR, if it exists. If it doesn't exist, does nothing",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := newPersistor()
			if err != nil {
				return err
			}
			runner := tipam.NewRunner(p, nil)
			return runner.Get(cidr)
		},
	}

	cmd.Flags().StringVar(&cidr, "cidr", "", "the cidr of the claim to be removed")

	return cmd
}
