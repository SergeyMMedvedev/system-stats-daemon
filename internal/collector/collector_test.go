package collector_test

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	_ "time"

	"github.com/SergeyMMedvedev/system-stats-daemon/internal/collector"
	c "github.com/SergeyMMedvedev/system-stats-daemon/internal/config"
	"github.com/stretchr/testify/require"
)

func TestCollectStatsLinux(t *testing.T) {
	os := runtime.GOOS
	if os != "linux" {
		t.Skip(fmt.Printf("skip wmic test for %s", os))
	}
	cfg := c.Config{
		StatsParams: c.StatsParamsConf{
			OS:          c.OSLinux,
			M:           2,
			N:           2,
			CPU:         true,
			DisksUsage:  true,
			DisksIoStat: true,
			NetStat:     true,
		},
	}

	statsCollector := collector.NewCollector(cfg)
	statsCollector.CollectStats(context.Background())
	var cpuStats [4]float64
	var disksStats []collector.DiskFree
	var disksIoStats []collector.DiskIoStat
	var netStats []collector.NetStat
	var tcpStats collector.TCPStats

	for i := 0; i < 3; i++ {
		cpuStats = <-statsCollector.CPUChan
		fmt.Printf("%+v", cpuStats)
		require.NotEqual(t, len(cpuStats), 0)
		disksStats = <-statsCollector.DisksFreeChan
		fmt.Printf("%+v", disksStats)
		disksIoStats = <-statsCollector.DisksIoStatChan
		fmt.Printf("%+v", disksIoStats)
		netStats = <-statsCollector.NetStatChan
		fmt.Printf("%+v", netStats)
		tcpStats = <-statsCollector.TCPStatChan
		fmt.Printf("%+v", tcpStats)
	}

	require.NotEqual(t, len(cpuStats), 0)
	require.NotEqual(t, len(disksStats), 0)
	require.NotEqual(t, len(disksIoStats), 0)
	require.NotNil(t, netStats, 0)
	require.NotNil(t, tcpStats, 0)
}

func TestCollectStatsWin(t *testing.T) {
	os := runtime.GOOS
	if os != "windows" {
		t.Skip(fmt.Printf("skip wmic test for %s", os))
	}
	cfg := c.Config{
		StatsParams: c.StatsParamsConf{
			OS:          c.OSWindows,
			M:           2,
			N:           2,
			CPU:         true,
			DisksUsage:  true,
			DisksIoStat: false,
			NetStat:     false,
		},
	}

	statsCollector := collector.NewCollector(cfg)
	statsCollector.CollectStats(context.Background())
	var cpuStats [4]float64
	var disksStats []collector.DiskFree

	for i := 0; i < 3; i++ {
		cpuStats = <-statsCollector.CPUChan
		fmt.Printf("%+v", cpuStats)
		require.NotEqual(t, len(cpuStats), 0)
		disksStats = <-statsCollector.DisksFreeChan
		fmt.Printf("%+v", disksStats)
	}

	require.NotEqual(t, len(cpuStats), 0)
	require.NotEqual(t, len(disksStats), 0)
}
