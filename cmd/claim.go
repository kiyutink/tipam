package cmd

import (
	"github.com/kiyutink/tipam/tipam"
	"github.com/spf13/cobra"
)

type claimFlags struct {
	complySubs bool
}

func newClaimCmd() *cobra.Command {
	var cidr string
	var tags []string
	var claimF claimFlags

	cmd := &cobra.Command{
		Use:   "claim",
		Short: "Create a block claim",
		Long:  "Create a block claim with provided CIDR and tags. The tags provided have to satisfy (read: include all the tags of) all the claims that contain this new claim.",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := newPersistor()
			if err != nil {
				return err
			}
			runner := tipam.NewRunner(p, nil)
			opts := []tipam.ClaimOption{}
			if claimF.complySubs {
				opts = append(opts, tipam.WithComplySubs(true))
			}
			return runner.Claim(cidr, tags, opts...)
		},
	}

	cmd.Flags().StringVar(&cidr, "cidr", "", "The CIDR range of the claim to be created")
	cmd.Flags().StringSliceVar(&tags, "tag", []string{}, "The list of tags to attach to a claim")
	cmd.Flags().BoolVar(&claimF.complySubs, "comply-subs", false, "Pass this flag to make subclaims comply with this claim by prepending tags")

	return cmd
}
