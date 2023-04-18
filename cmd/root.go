package cmd

import (
	"github.com/kiyutink/tipam/internal/visual"
	"github.com/kiyutink/tipam/tipam"
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
			runner := tipam.NewRunner(p, nil)
			visual.InitTipam(runner)
			return nil
		},
	}

	// Flags used by the persistor
	rootCmd.PersistentFlags().StringVar(&persistF.persistorType, "persist", "localyaml", "Which persistor to use. Only 'localyaml' is available at the time, which is also the default value")
	rootCmd.PersistentFlags().StringVar(&persistF.localYAMLFileName, "filename", "tipam.yaml", "The filename for the 'localyaml' persistor. Defaults to 'tipam.yaml'")

	// Flags mapped into runner options
	rootCmd.PersistentFlags().BoolVar(&runnerOpts.DoLock, "do-lock", false, "Whether or not use state locking. Defaults to false")

	rootCmd.AddCommand(newClaimCmd())
	rootCmd.AddCommand(newReleaseCmd())
	rootCmd.AddCommand(newGetCmd())

	return rootCmd
}
