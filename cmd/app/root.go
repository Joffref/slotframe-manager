package app

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "slotframe-manager",
		Short: "slotframe-manager is a CLI tool for managing slotframes on 6TiSCH networks",
		Long:  "slotframe-manager is a CLI tool for managing slotframes on 6TiSCH networks",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
