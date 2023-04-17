package cmd

import (
	"github.com/kiyutink/tipam/core"
	"github.com/kiyutink/tipam/persist"
	"github.com/spf13/cobra"
)

func newReserveCmd() *cobra.Command {
	var cidr string
	var tags []string
	var reserveFlags core.ReserveFlags

	cmd := &cobra.Command{
		Use:   "reserve",
		Short: "Create a block reseration",
		Long:  "Create a block reservation with provided CIDR and tags. The tags provided have to satisfy (read: include all the tags of) all the reservations that contain this new reservation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := persist.NewPersistor(persistFlags)
			if err != nil {
				return err
			}
			runner := core.NewRunner(p)
			return runner.Reserve(cidr, tags, reserveFlags)
		},
	}

	cmd.Flags().StringVar(&cidr, "cidr", "", "The CIDR range of the reservation to be created")
	cmd.Flags().StringSliceVar(&tags, "tag", []string{}, "The list of tags to attach to a reservation")
	cmd.Flags().BoolVar(&reserveFlags.ComplySubs, "comply-subs", false, "Pass this flag to make subreservations comply with this reservation by prepending tags")

	return cmd
}
