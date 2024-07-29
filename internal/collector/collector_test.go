package collector_test

import (
	"context"
	"fmt"
	"testing"
	_ "time"

	"github.com/SergeyMMedvedev/system-stats-daemon/internal/collector"
	c "github.com/SergeyMMedvedev/system-stats-daemon/internal/config"
	"github.com/stretchr/testify/require"
)

var cfg = c.Config{
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

func TestCollectStats(t *testing.T) {
	statsCollector := collector.NewCollector(cfg)
	statsCollector.CollectStats(context.Background())
	var cpuStats [4]float64
	var disksStats []collector.DiskFree
	var disksIoStats []collector.DiskIoStat
	var netStats []collector.NetStat
	var tcpStats collector.TCPStats

	for i := 0; i < 3; i++ {
		cpuStats = <-statsCollector.CpuChan
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
