package app

import "github.com/spf13/cobra"

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the slotframe-manager",
		Long:  "Run the slotframe-manager",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)
