package cpu_test

import (
	"fmt"
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/internal/collector/cpu"
	c "github.com/SergeyMMedvedev/system-stats-daemon/internal/config"
	"github.com/stretchr/testify/require"
)

var cfg = c.Config{
	StatsParams: c.StatsParamsConf{
		OS: c.OSLinux,
	},
}

func TestCollectCPUStats(t *testing.T) {
	loadAvg, us, sy, id, err := cpu.CollectCPUStats(cfg.StatsParams.OS)
	fmt.Printf("loadAvg: %f, us: %f, sy: %f, id: %f", loadAvg, us, sy, id)
	require.NoError(t, err)
}
