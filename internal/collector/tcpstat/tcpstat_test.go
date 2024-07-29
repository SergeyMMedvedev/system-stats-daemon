package tcpstat_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/internal/collector/tcpstat"
	"github.com/stretchr/testify/require"
)

func TestCollectTCPStats(t *testing.T) {
	os := runtime.GOOS
	if os != "linux" {
		t.Skip(fmt.Printf("skip wmic test for %s", os))
	}
	_, err := tcpstat.CollectTCPStats()
	require.NoError(t, err)
}
