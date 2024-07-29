package iostat_test

import (
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/internal/collector/iostat"
	"github.com/stretchr/testify/require"
)

func TestCollectIoStat(t *testing.T) {
	iostats, err := iostat.CollectIoStat()
	require.NoError(t, err)
	require.NotEqual(t, len(iostats), 0)
}
