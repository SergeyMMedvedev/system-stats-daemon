package df_test

import (
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/pkg/df"
	"github.com/stretchr/testify/require"
)

func TestDiskFreeInodeStats(t *testing.T) {
	diskInfo, err := df.DiskFreeInodeStats()
	require.NoError(t, err)
	require.Contains(t, diskInfo, "Inodes")
}

func TestDiskFreeStats(t *testing.T) {
	diskInfo, err := df.DiskFreeStats()
	require.NoError(t, err)
	require.Contains(t, diskInfo, "Available")
}
