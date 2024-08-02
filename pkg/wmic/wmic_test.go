package wmic_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/pkg/wmic"
	"github.com/stretchr/testify/require"
)

func TestCPUGetLoadPercentage(t *testing.T) {
	os := runtime.GOOS
	if os != "windows" {
		t.Skip(fmt.Printf("skip wmic test for %s", os))
	}
	info, err := wmic.CPUGetLoadPercentage()
	require.NoError(t, err)
	fmt.Println(info)
}
