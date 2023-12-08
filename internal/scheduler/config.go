package scheduler

import (
	"fmt"
	"github.com/Joffref/slotframe-manager/internal/api"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	FrameSize         int           `mapstructure:"frameSize"`
	NumberCh          int           `mapstructure:"numberCh"`
	LogLevel          string        `mapstructure:"logLevel"`
	KeepAliveInterval time.Duration `mapstructure:"keepAliveInterval"`
	KeepAliveTimeout  time.Duration `mapstructure:"keepAliveTimeout"`
	api.APIConfig     `mapstructure:",squash"`
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

	if c.KeepAliveInterval == 0 {
		c.KeepAliveInterval = 10 * time.Second
	}

	if c.KeepAliveTimeout == 0 {
		c.KeepAliveTimeout = 30 * time.Second
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
