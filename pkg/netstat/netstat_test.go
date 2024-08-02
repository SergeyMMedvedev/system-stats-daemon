package netstat_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/pkg/netstat"
	"github.com/stretchr/testify/require"
)

func TestGetNetstat(t *testing.T) {
	os := runtime.GOOS
	if os != "linux" {
		t.Skip(fmt.Printf("skip wmic test for %s", os))
	}
	diskInfo, err := netstat.NetStat()
	require.NoError(t, err)
	require.Contains(t, diskInfo, "Active Internet connections")
}
