package netstat_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/internal/collector/netstat"
	"github.com/stretchr/testify/require"
)

func TestCollectNetstat(t *testing.T) {
	os := runtime.GOOS
	if os != "linux" {
		t.Skip(fmt.Printf("skip wmic test for %s", os))
	}
	_, err := netstat.CollectNetstat()
	require.NoError(t, err)
}
