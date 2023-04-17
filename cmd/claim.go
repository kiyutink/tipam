package cmd

import (
	"github.com/kiyutink/tipam/core"
	"github.com/spf13/cobra"
)

func newClaimCmd() *cobra.Command {
	var cidr string
	var tags []string
	var claimFlags core.ClaimFlags

	cmd := &cobra.Command{
		Use:   "claim",
		Short: "Create a block claim",
		Long:  "Create a block claim with provided CIDR and tags. The tags provided have to satisfy (read: include all the tags of) all the claims that contain this new claim.",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := newPersistor()
			if err != nil {
				return err
			}
			runner := core.NewRunner(p, nil)
			return runner.Claim(cidr, tags, claimFlags)
		},
	}

	cmd.Flags().StringVar(&cidr, "cidr", "", "The CIDR range of the claim to be created")
	cmd.Flags().StringSliceVar(&tags, "tag", []string{}, "The list of tags to attach to a claim")
	cmd.Flags().BoolVar(&claimFlags.ComplySubs, "comply-subs", false, "Pass this flag to make subclaims comply with this claim by prepending tags")

	return cmd
}
