package ss_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/pkg/ss"
	"github.com/stretchr/testify/require"
)

func TestGrepTCP(t *testing.T) {
	os := runtime.GOOS
	if os != "linux" {
		t.Skip(fmt.Printf("skip wmic test for %s", os))
	}
	info, err := ss.GrepTCP()
	require.NoError(t, err)
	require.Contains(t, info, "TCP")
}
