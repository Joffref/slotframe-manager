package app

import (
	"github.com/Joffref/slotframe-manager/internal/scheduler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configPath string
	cfg        scheduler.Config
	parentId   uint16
	id         uint16
	etx        uint16
	rootCmd    = &cobra.Command{
		Use:   "slotframe-manager-testclient",
		Short: "slotframe-manager-testclient is a CLI tool for testing slotframes on RPL networks",
		Long:  "slotframe-manager-testclient is a CLI tool for testing slotframes on RPL networks",
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to the configuration file")
	runCmd.Flags().Uint16VarP(&parentId, "parentId", "p", 0, "Parent ID")
	runCmd.Flags().Uint16VarP(&id, "id", "i", 0, "ID")
	runCmd.Flags().Uint16VarP(&etx, "etx", "e", 0, "ETX")
	viper.BindPFlag("parentId", runCmd.Flags().Lookup("parentId"))
	viper.BindPFlag("id", runCmd.Flags().Lookup("id"))
	viper.BindPFlag("etx", runCmd.Flags().Lookup("etx"))
}

func readInConfig() error {
	if configPath == "" {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
	}
	viper.SetConfigFile(configPath)
	return viper.ReadInConfig()
}

func parseConfig() error {
	return viper.Unmarshal(&cfg)
}

// Execute executes the root command.
func Execute() error {
	if err := readInConfig(); err != nil {
		return err
	}
	if err := parseConfig(); err != nil {
		return err
	}
	return rootCmd.Execute()
}
