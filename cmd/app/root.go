package app

import (
	"github.com/Joffref/slotframe-manager/internal/scheduler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configPath string
	cfg        scheduler.Config
	rootCmd    = &cobra.Command{
		Use:   "slotframe-manager",
		Short: "slotframe-manager is a CLI tool for managing slotframes on RPL networks",
		Long:  "slotframe-manager is a CLI tool for managing slotframes on RPL networks",
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to the configuration file")
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
