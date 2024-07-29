package iostat_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/internal/collector/iostat"
	"github.com/stretchr/testify/require"
)

func TestCollectIoStat(t *testing.T) {
	os := runtime.GOOS
	if os != "linux" {
		t.Skip(fmt.Printf("skip wmic test for %s", os))
	}
	iostats, err := iostat.CollectIoStat()
	require.NoError(t, err)
	require.NotEqual(t, len(iostats), 0)
}
