package iostat_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/pkg/iostat"
	"github.com/stretchr/testify/require"
)

func TestGetIOStat(t *testing.T) {
	os := runtime.GOOS
	if os != "linux" {
		t.Skip(fmt.Printf("skip wmic test for %s", os))
	}
	diskInfo, err := iostat.GetStat()
	require.NoError(t, err)
	require.Contains(t, string(diskInfo), "sysstat")
}
