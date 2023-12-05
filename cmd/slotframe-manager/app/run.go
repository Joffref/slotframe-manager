package app

import (
	"github.com/Joffref/slotframe-manager/internal/api"
	"github.com/Joffref/slotframe-manager/internal/graph"
	"github.com/Joffref/slotframe-manager/internal/scheduler"
	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the slotframe-manager",
		Long:  "Run the slotframe-manager",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cfg.Validate(); err != nil {
				return err
			}
			scheduler := scheduler.NewScheduler(graph.NewDoDAG(), &cfg)
			go scheduler.Schedule()
			go scheduler.ControlLoop()
			api, err := api.NewAPI(&cfg.APIConfig, scheduler)
			if err != nil {
				return err
			}
			return api.Run()
		},
	}
)
