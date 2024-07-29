package netstat_test

import (
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/pkg/netstat"
	"github.com/stretchr/testify/require"
)

func TestGetnetstat(t *testing.T) {
	diskInfo, err := netstat.NetStat()
	require.NoError(t, err)
	require.Contains(t, diskInfo, "Active Internet connections")
}
