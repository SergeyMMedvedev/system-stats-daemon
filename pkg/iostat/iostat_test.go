package iostat_test

import (
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/pkg/iostat"
	"github.com/stretchr/testify/require"
)

func TestGetIOStat(t *testing.T) {
	diskInfo, err := iostat.GetStat()
	require.NoError(t, err)
	require.Contains(t, string(diskInfo), "sysstat")
}
