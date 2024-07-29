package tcpstat_test

import (
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/internal/collector/tcpstat"
	"github.com/stretchr/testify/require"
)

func TestCollectTCPStats(t *testing.T) {
	_, err := tcpstat.CollectTCPStats()
	require.NoError(t, err)
}
