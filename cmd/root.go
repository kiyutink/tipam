package cmd

import (
	"github.com/kiyutink/tipam/core"
	"github.com/kiyutink/tipam/persist"
	"github.com/kiyutink/tipam/tipam"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tipam",
		Short: "tipam is an IP Address Manager for the terminal",
		Run: func(cmd *cobra.Command, args []string) {
			tipam.InitTipam()
		},
	}

	yamlReservationsClient := &persist.YamlReservationsClient{}
	runner := &core.Runner{
		ReservationsClient: yamlReservationsClient,
	}
	rootCmd.AddCommand(newReserveCmd(runner))
	rootCmd.AddCommand(newReleaseCmd(runner))

	return rootCmd
}
