package cmd

import (
	"fmt"

	"github.com/kiyutink/tipam/persist"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Show the list of claimed CIDRs",
		Long:  "Show the list of claimed CIDRs. Currenly only shows the enitre state",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := newPersistor()
			if err != nil {
				return err
			}
			runner := newRunner(p)

			state, err := runner.List()
			if err != nil {
				return err
			}

			yamlState, err := persist.StateToYAMLString(state)
			if err != nil {
				return err
			}
			fmt.Println(yamlState)
			return nil
		},
	}

	return cmd
}
