package cpu_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/internal/collector/cpu"
	c "github.com/SergeyMMedvedev/system-stats-daemon/internal/config"
	"github.com/stretchr/testify/require"
)

func TestCollectCPUStatsLinux(t *testing.T) {
	os := runtime.GOOS
	if os != "linux" {
		t.Skip(fmt.Printf("skip wmic test for %s", os))
	}
	cfg := c.Config{
		StatsParams: c.StatsParamsConf{
			OS: c.OSLinux,
		},
	}
	cpuStats, err := cpu.CollectCPUStats(cfg.StatsParams.OS)
	fmt.Printf("loadAvg: %f, us: %f, sy: %f, id: %f", cpuStats.L, us, sy, id)
	require.NoError(t, err)
}

func TestCollectCPUStatsWin(t *testing.T) {
	os := runtime.GOOS
	if os != "windows" {
		t.Skip(fmt.Printf("skip wmic test for %s", os))
	}
	cfg := c.Config{
		StatsParams: c.StatsParamsConf{
			OS: c.OSWindows,
		},
	}
	loadAvg, us, sy, id, err := cpu.CollectCPUStats(cfg.StatsParams.OS)
	fmt.Printf("loadAvg: %f, us: %f, sy: %f, id: %f", loadAvg, us, sy, id)
	require.NoError(t, err)
}
