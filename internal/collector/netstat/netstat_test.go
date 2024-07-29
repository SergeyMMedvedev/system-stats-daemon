package netstat_test

import (
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/internal/collector/netstat"
	"github.com/stretchr/testify/require"
)

func TestCollectNetstat(t *testing.T) {
	_, err := netstat.CollectNetstat()
	require.NoError(t, err)
}
