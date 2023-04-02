package cmd

import (
	"github.com/spf13/cobra"
)

type CreateReservationRunner interface {
	CreateReservation(cidr string, tags []string) error
}

func newReserveCmd(runner CreateReservationRunner) *cobra.Command {
	var cidr string
	var tags []string

	cmd := &cobra.Command{
		Use:   "reserve",
		Short: "Creates a block reseration",
		Long:  "Creates a block reservation with provided CIDR and tags. The tags provided have to satisfy (read: include all the tags of) all the reservations that contain this new reservation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner.CreateReservation(cidr, tags)
		},
	}

	cmd.Flags().StringVar(&cidr, "cidr", "", "The CIDR range of the reservation to be created")
	cmd.Flags().StringSliceVar(&tags, "tag", []string{}, "The list of tags to attach to a reservation")

	return cmd
}
