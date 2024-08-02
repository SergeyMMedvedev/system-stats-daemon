package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	c "github.com/SergeyMMedvedev/system-stats-daemon/internal/config"
	"github.com/SergeyMMedvedev/system-stats-daemon/internal/server"
)

var (
	configFile string
	port       int
	host       string
)

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.yaml", "Path to configuration file")
	flag.StringVar(&host, "host", "localhost", "Host")
	flag.IntVar(&port, "port", 50051, "Server port")
}

func main() {
	flag.Parse()
	fmt.Println("host", host)
	fmt.Println("port", port)
	config := c.NewConfig()
	err := config.Read(configFile)
	if err != nil {
		fmt.Printf("failed to read config: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", config)

	grpcServer := server.NewServer(config, host, port)

	ctx := context.Background()
	err = grpcServer.Run(ctx)
	if err != nil {
		fmt.Printf("failed to run server: %v\n", err)
		os.Exit(1)
	}
}
