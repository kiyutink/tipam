package cmd

import (
	"github.com/kiyutink/tipam/internal/visual"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tipam",
		Short: "tipam is an IP Address Manager for the terminal",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := newPersistor()
			if err != nil {
				return err
			}
			runner := newRunner(p)
			visual.InitTipam(runner)
			return nil
		},
	}

	// Flags used by the persistor
	rootCmd.PersistentFlags().StringVar(&persistF.persistor, "persistor", "localyaml", "which persistor to use. Only 'localyaml' is available at the time")
	rootCmd.PersistentFlags().StringVar(&persistF.localYAMLFileName, "filename", "tipam.yaml", "the filename for the 'localyaml' persistor")

	// Flags mapped into runner options
	rootCmd.PersistentFlags().BoolVar(&runnerF.lock, "lock", false, "use state locking (for concurrent runs)")

	rootCmd.AddCommand(newClaimCmd())
	rootCmd.AddCommand(newReleaseCmd())
	rootCmd.AddCommand(newGetCmd())

	return rootCmd
}
