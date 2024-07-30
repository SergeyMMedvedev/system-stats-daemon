package integration_test

import (
	"context"
	"fmt"
	_ "os"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/SergeyMMedvedev/system-stats-daemon/internal/client"
	"github.com/SergeyMMedvedev/system-stats-daemon/internal/config"
	"github.com/SergeyMMedvedev/system-stats-daemon/internal/pb"
	"github.com/stretchr/testify/require"
	_ "google.golang.org/grpc"
)

func TestIntegration(t *testing.T) {
	// Step 1: Start Docker Compose
	cmd := exec.Command("docker", "compose", "up", "--build", "-d")
	fmt.Println("Starting Docker Compose...")
	err := cmd.Run()
	require.NoError(t, err)

	// Step 2: Defer Docker Compose down
	defer func() {
		cmd := exec.Command("docker", "compose", "down")
		err := cmd.Run()
		require.NoError(t, err)
	}()
	fmt.Println("Docker Compose started successfully.")
	// Step 3: Wait for the server to be ready
	time.Sleep(5 * time.Second) // Adjust the wait time as needed
	c := client.NewClient()
	configFile := "integration_config.yaml"
	cfg := config.NewConfig()
	err = cfg.Read(configFile)
	fmt.Printf("Config read successfully: %+v\n", cfg)
	require.NoError(t, err)
	err = c.Connect("localhost", 50051, cfg)
	require.NoError(t, err)
	fmt.Println("Connected to server")
	require.NoError(t, err)
	ctx := context.Background()
	timeout := time.Second * time.Duration(cfg.StatsParams.N*40)
	ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		fmt.Println("Start stream stats.")
		err = c.StreamSystemStats(ctxTimeout, &pb.SystemStatsRequest{})
		if err != nil {
			require.ErrorIs(t, err, context.DeadlineExceeded)
			fmt.Println("Stop collecting system stats: deadline exceeded.")
		}
		wg.Done()
	}()
	time.Sleep(time.Duration(cfg.StatsParams.N*20) * time.Second)
	// Step 4: Execute cpu stress test on the server
	fmt.Println("Executing CPU stress test on the server...")
	cmd = exec.Command("docker", "compose", "exec", "server", "stress", "-q", "-t", "40s", "--cpu", "16")
	err = cmd.Run()
	require.NoError(t, err)

	wg.Wait()
	stats := c.SystemStatsResponses
	for _, stat := range stats {
		fmt.Printf("CPU LoadAverage: %f\n", stat.LoadAverage)
		fmt.Printf("CPU Stats: %+v\n", stat.CPUStats)
	}
}
