package top_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/pkg/top"
	"github.com/stretchr/testify/require"
)

func TestTop(t *testing.T) {
	os := runtime.GOOS
	if os != "linux" {
		t.Skip(fmt.Printf("skip wmic test for %s", os))
	}
	info, err := top.Top()
	require.NoError(t, err)
	require.Contains(t, info, "load average")
}
