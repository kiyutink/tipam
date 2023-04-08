package cmd

import (
	"github.com/kiyutink/tipam/core"
	"github.com/kiyutink/tipam/persist"
	"github.com/kiyutink/tipam/tipam"
	"github.com/spf13/cobra"
)

var persistFlags = persist.PersistFlags{}

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tipam",
		Short: "tipam is an IP Address Manager for the terminal",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := persist.NewPersistor(persistFlags)
			if err != nil {
				return err
			}
			runner := core.NewRunner(p)
			tipam.InitTipam(runner)
			return nil
		},
	}
	rootCmd.PersistentFlags().StringVar(&persistFlags.PersistorType, "persist", "localyaml", "Which persistor to use. Only 'localyaml' is available at the time, which is also the default value")
	rootCmd.PersistentFlags().StringVar(&persistFlags.LocalYAMLFileName, "filename", "tipam.yaml", "The filename for the 'localyaml' persistor. Defaults to 'tipam.yaml'")

	rootCmd.AddCommand(newReserveCmd())
	rootCmd.AddCommand(newReleaseCmd())

	return rootCmd
}
