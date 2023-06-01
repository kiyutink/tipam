package cmd

import (
	"github.com/kiyutink/tipam/tipam"
	"github.com/spf13/cobra"
)

func newClaimCmd() *cobra.Command {
	type claimFlags struct {
		// Required flags
		cidr string
		tags []string

		// Optional flags
		complySubs bool
		final      bool
	}

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

			runner := newRunner(p)

			return runner.Claim(tipam.MustParseClaimFromCIDR(claimF.cidr, claimF.tags, claimF.final), tipam.ClaimOpts{ComplySubs: claimF.complySubs})
		},
	}

	cmd.Flags().StringVar(&claimF.cidr, "cidr", "", "the CIDR range of the claim to be created")
	cmd.MarkFlagRequired("cidr")

	cmd.Flags().StringSliceVar(&claimF.tags, "tag", []string{}, "the list of tags attach to created claim; multiple instances of this flag can be passed to attach multiple flags")
	cmd.MarkFlagRequired("tag")

	cmd.Flags().BoolVar(&claimF.complySubs, "comply-subs", false, "pass this flag to make subclaims comply with this claim by prepending tags")
	cmd.Flags().BoolVar(&claimF.final, "final", false, "pass this flag to make the claim final (disallow subclaims)")

	return cmd
}
