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

	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false

	// Determines the persistor to use
	rootCmd.PersistentFlags().StringVar(&persistF.persistor, "persistor", "localyaml", "which persistor to use. Available options: localyaml, inmemory")

	// For localyaml persistor
	rootCmd.PersistentFlags().StringVar(&persistF.localYAMLFileName, "filename", "tipam.yaml", "'localyaml' persistor - filename to use")

	// For s3dynamo persistor
	rootCmd.PersistentFlags().StringVar(&persistF.s3DynamoBucket, "bucket", "", "'s3dynamo' persistor - the s3 bucket to use")
	rootCmd.PersistentFlags().StringVar(&persistF.s3DynamoKeyInBucket, "key-in-bucket", "tipam.yaml", "'s3dynamo' persistor - the key to use in the s3 bucket")
	rootCmd.PersistentFlags().StringVar(&persistF.s3DynamoTable, "table", "", "'s3dynamo' persistor - the DynamoDB table to use")
	rootCmd.PersistentFlags().IntVar(&persistF.s3DynamoLeaseDuration, "lease-duration", 10, "'s3dynamo' persistor - how long to hold the lock (in seconds)")
	rootCmd.PersistentFlags().IntVar(&persistF.s3DynamoPollInterval, "poll-interval", 3, "'s3dynamo' persistor - how often to poll dynamodb to check whether the lock has been released")

	// Flags mapped into runner options
	rootCmd.PersistentFlags().BoolVar(&runnerF.lock, "lock", true, "use state locking (for concurrent runs)")

	rootCmd.AddCommand(newClaimCmd())
	rootCmd.AddCommand(newReleaseCmd())
	rootCmd.AddCommand(newGetCmd())

	return rootCmd
}
