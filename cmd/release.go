package cmd

import (
	"github.com/spf13/cobra"
)

func newReleaseCmd() *cobra.Command {
	var cidr string
	cmd := &cobra.Command{
		Use:   "release",
		Short: "release a block claim",
		Long:  "release the block claim with the provided cidr",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := newPersistor()
			if err != nil {
				return err
			}
			runner := newRunner(p)
			return runner.Release(cidr)
		},
	}

	cmd.Flags().StringVar(&cidr, "cidr", "", "The CIDR range of the block to be released")
	return cmd
}
