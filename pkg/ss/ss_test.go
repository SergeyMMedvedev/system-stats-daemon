package ss_test

import (
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/pkg/ss"
	"github.com/stretchr/testify/require"
)

func TestGrepTCP(t *testing.T) {
	info, err := ss.GrepTCP()
	require.NoError(t, err)
	require.Contains(t, info, "TCP")
}
