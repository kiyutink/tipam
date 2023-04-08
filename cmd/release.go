package cmd

import (
	"github.com/kiyutink/tipam/core"
	"github.com/kiyutink/tipam/persist"
	"github.com/spf13/cobra"
)

func newReleaseCmd() *cobra.Command {
	var cidr string
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Release a block reservation or claim",
		Long:  "Release the block reservation or claim with the provided cidr",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := persist.NewPersistor(persistFlags)
			if err != nil {
				return err
			}
			runner := core.NewRunner(p)
			return runner.Release(cidr)
		},
	}

	cmd.Flags().StringVar(&cidr, "cidr", "", "The CIDR range of the block to be released")
	return cmd
}
