package ringbuffer_test

import (
	"testing"

	"github.com/SergeyMMedvedev/system-stats-daemon/internal/ringbuffer"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	rb := ringbuffer.NewRingBuffer(3)
	require.True(t, rb.IsEmpty())
	rb.Enqueue(1)
	rb.Enqueue(2)
	rb.Enqueue(3)

	require.False(t, rb.IsEmpty())
	require.True(t, rb.IsFull())
	require.Equal(t, rb.Size(), 3)
	require.Equal(t, rb.Average(), float64(2))

	rb.Enqueue(10)
	require.True(t, rb.IsFull())
	require.Equal(t, rb.Size(), 3)
	require.Equal(t, rb.Average(), float64(5))
}
