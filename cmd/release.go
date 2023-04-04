package cmd

import "github.com/spf13/cobra"

type Releaser interface {
	Release(cidr string) error
}

func newReleaseCmd(r Releaser) *cobra.Command {
	var cidr string
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Release a block reservation or claim",
		Long:  "Release the block reservation or claim with the provided cidr",
		RunE: func(cmd *cobra.Command, args []string) error {
			return r.Release(cidr)
		},
	}

	cmd.Flags().StringVar(&cidr, "cidr", "", "The CIDR range of the block to be released")
	return cmd
}
