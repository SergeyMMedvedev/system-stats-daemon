package integration_test

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/SergeyMMedvedev/system-stats-daemon/internal/client"
	"github.com/SergeyMMedvedev/system-stats-daemon/internal/config"
	"github.com/SergeyMMedvedev/system-stats-daemon/internal/pb"
	"github.com/stretchr/testify/require"
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

	// prepare stats collecting params
	collectingIters := 10
	statsCollectingTimeout := time.Second * time.Duration(
		cfg.StatsParams.N*collectingIters,
	)
	halfPeriod := statsCollectingTimeout / 2
	fmt.Println("Stats collecting timeout:", statsCollectingTimeout)
	fmt.Println("halfPeriod:", halfPeriod)
	cpuStressPeriod := fmt.Sprintf("%ds", cfg.StatsParams.N*collectingIters/2)
	cpuLoadThreads := "16"
	fmt.Println("CPU stress period:", cpuStressPeriod)
	fmt.Println("CPU load threads:", cpuLoadThreads)

	// Step 4: Collect system stats
	ctxTimeout, cancel := context.WithTimeout(ctx, statsCollectingTimeout)
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
	fmt.Println(
		"Collect statistics for the first half of the time without loading the processor...",
	)
	time.Sleep(halfPeriod)
	// Step 4: Execute cpu stress test on the server
	fmt.Println("Executing CPU stress test on the server...")
	cmd = exec.Command(
		"docker", "compose", "exec", "server", "stress", "-q", "-t", cpuStressPeriod, "--cpu", cpuLoadThreads,
	)
	err = cmd.Run()
	require.NoError(t, err)

	wg.Wait()
	stats := c.SystemStatsResponses
	statsWithoutStress := stats[:len(stats)/2]
	statsWithStress := stats[len(stats)/2:]
	lastStat := stats[len(stats)-1]

	getAverageUserModeCPULoad := func(stats []*client.Stats) float32 {
		sum := float32(0)
		for _, stat := range stats {
			sum += float32(stat.CPUStats.UserMode)
		}
		return sum / float32(len(stats))
	}

	// Step 5: Verify system stats
	averageUserModeCPULoadWithoutStress := getAverageUserModeCPULoad(statsWithoutStress)
	averageUserModeCPULoadWithStress := getAverageUserModeCPULoad(statsWithStress)

	fmt.Println("Average User Mode CPU Load without stress:", averageUserModeCPULoadWithoutStress)
	fmt.Println("Average User Mode CPU Load with stress:", averageUserModeCPULoadWithStress)
	require.True(t, averageUserModeCPULoadWithStress > averageUserModeCPULoadWithoutStress)
	require.True(t, len(lastStat.DisksFree) > 0)
	require.True(t, len(lastStat.DisksIoStat) > 0)
	require.True(t, len(lastStat.NetStat) > 0)
	estabTCPconns := lastStat.TCPStats.Estab
	require.True(t, estabTCPconns > 0)
	fmt.Println("Integration test passed successfully.")
}
