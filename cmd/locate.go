package cmd

import (
	"fmt"
	"strings"

	"github.com/kiyutink/tipam/tipam"
	"github.com/spf13/cobra"
)

func newLocateCmd() *cobra.Command {
	type locateFlags struct {
		// Required flags
		within  string
		maskLen int

		// Optional flags
		claim bool
		tags  string
		final bool
	}

	locateF := locateFlags{}

	cmd := &cobra.Command{
		Use:   "locate",
		Short: "Locate a new claim",
		Long:  "Locate a subnet. The new subnet will be created in the next available space where the claim would fit, within a parent block with given tags",

		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := newPersistor()
			if err != nil {
				return err
			}

			runner := newRunner(p)
			withinSlice := strings.Split(locateF.within, "/")

			cidr, err := runner.Locate(withinSlice, locateF.maskLen)
			if err != nil {
				return err
			}

			fmt.Println(cidr)

			if locateF.claim {
				tagsSlice := strings.Split(locateF.tags, "/")
				claim := tipam.MustParseClaimFromCIDR(cidr, tagsSlice, locateF.final)

				return runner.Claim(claim, tipam.ClaimOpts{})
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&locateF.within, "within", "", "The flags used to locate claims, within which to locate the new claim")
	cmd.MarkFlagRequired("within")

	cmd.Flags().IntVar(&locateF.maskLen, "mask-len", 0, "The length of the mask of the claimed subnet")
	cmd.MarkFlagRequired("mask-len")

	cmd.Flags().BoolVar(&locateF.claim, "claim", false, "Whether to immediately claim the located subnet")
	cmd.Flags().StringVar(&locateF.tags, "tags", "", "The tags to tag the new claim with (only relevant if --claim is set to true)")
	cmd.MarkFlagsRequiredTogether("claim", "tags")
	cmd.Flags().BoolVar(&locateF.final, "final", false, "Whether the new claim is final, defaults to false (only relevant if --claim is set to true)")

	return cmd
}
