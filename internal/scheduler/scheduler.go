package scheduler

import (
	"fmt"
	"github.com/Joffref/slotframe-manager/internal/api"
	"github.com/Joffref/slotframe-manager/internal/graph"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	FrameSize     int    `mapstructure:"frameSize"`
	NumberCh      int    `mapstructure:"numberCh"`
	LogLevel      string `mapstructure:"logLevel"`
	api.APIConfig `mapstructure:",squash"`
}

func (c *Config) Validate() error {
	if c.FrameSize <= 0 {
		return fmt.Errorf("invalid frame size: %d, must be greater than 0", c.FrameSize)
	}
	if c.NumberCh <= 0 {
		return fmt.Errorf("invalid number of channels: %d, must be greater than 0", c.NumberCh)
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	if c.Address == "" {
		c.Address = ":5688"
	}
	var level slog.Level
	switch c.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	default:
		slog.Error("Invalid log level, using info")
		level = slog.LevelInfo
	}
	opts := slog.HandlerOptions{
		Level: level,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &opts))
	slog.SetDefault(logger)
	return nil
}

type Scheduler struct {
	cfg   *Config
	Frame *Frame
	dodag *graph.DoDAG
}

func NewScheduler(dodag *graph.DoDAG, cfg *Config) *Scheduler {
	slog.Debug("Creating scheduler")
	return &Scheduler{
		cfg:   cfg,
		Frame: NewFrame(cfg.FrameSize, cfg.NumberCh),
		dodag: dodag,
	}
}

func (s *Scheduler) Register(parentId string, id string, etx int, input api.Slots) (*api.Slots, error) {
	if parentId == "0" {
		parentId = ""
	}
	if !s.dodag.IsNode(parentId) && parentId != "" {
		return nil, api.ErrorParentNodeDoesNotExist{ParentId: parentId}
	}
	if s.dodag.IsNode(id) {
		s.dodag.Nodes[id].LastSeen = time.Now()
		return &api.Slots{
			EmittingSlots:  s.dodag.Nodes[id].EmittingSlots,
			ListeningSlots: s.dodag.Nodes[id].ListeningSlots,
		}, nil
	}
	node := graph.NewNode(parentId, id, etx, input.EmittingSlots, input.ListeningSlots)
	s.dodag.AddNode(node)
	return &api.Slots{
		EmittingSlots:  node.EmittingSlots,
		ListeningSlots: node.ListeningSlots,
	}, nil
}

func (s *Scheduler) Schedule() {
	for {
		currentVersion := s.Frame.Version
		frame, err := ComputeFrame(s.dodag, s.cfg.FrameSize, s.cfg.NumberCh)
		if err != nil {
			slog.Error(fmt.Sprintf("cannot compute frame: %v", err))
			time.Sleep(5 * time.Second)
			continue
		}
		s.Frame = frame
		s.Frame.Version = currentVersion + 1
		time.Sleep(5 * time.Second)
	}
}

func (s *Scheduler) Version() int {
	return s.Frame.Version
}
