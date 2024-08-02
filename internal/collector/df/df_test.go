package df_test

import (
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/internal/collector/df"
	"github.com/stretchr/testify/require"
)

func TestCollectDiskFreeStats(t *testing.T) {
	diskInfo, err := df.CollectDiskFreeStats()
	require.NoError(t, err)
	require.NotEqual(t, len(diskInfo), 0)
}

func TestCollectDiskFreeInodeStats(t *testing.T) {
	diskInfo, err := df.CollectDiskFreeInodeStats()
	require.NoError(t, err)
	require.NotEqual(t, len(diskInfo), 0)
}
