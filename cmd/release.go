package cmd

import (
	"github.com/kiyutink/tipam/tipam"
	"github.com/spf13/cobra"
)

func newReleaseCmd() *cobra.Command {
	var cidr string
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Release a block claim or claim",
		Long:  "Release the block claim or claim with the provided cidr",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := newPersistor()
			if err != nil {
				return err
			}
			runner := tipam.NewRunner(p, nil)
			return runner.Release(cidr)
		},
	}

	cmd.Flags().StringVar(&cidr, "cidr", "", "The CIDR range of the block to be released")
	return cmd
}
