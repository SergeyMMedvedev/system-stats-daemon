package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/SergeyMMedvedev/system-stats-daemon/internal/client"
	"github.com/SergeyMMedvedev/system-stats-daemon/internal/config"
	"github.com/SergeyMMedvedev/system-stats-daemon/internal/pb"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.yaml", "Path to configuration file")
}

func main() {
	c := client.NewClient()

	cfg := config.NewConfig()
	err := cfg.Read(configFile)
	if err != nil {
		fmt.Printf("failed to read config: %v\n", err)
		os.Exit(1)
	}
	err = c.Connect("localhost", 50051, cfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	ctx := context.Background()
	err = c.StreamSystemStats(ctx, &pb.SystemStatsRequest{})
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
