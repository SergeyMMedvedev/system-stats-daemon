package top_test

import (
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/pkg/top"
	"github.com/stretchr/testify/require"
)

func TestTop(t *testing.T) {
	info, err := top.Top()
	require.NoError(t, err)
	require.Contains(t, info, "load average")
}
