package server_test

import (
	"context"
	"testing"
	"time"

	c "github.com/SergeyMMedvedev/system-stats-daemon/internal/config"
	"github.com/SergeyMMedvedev/system-stats-daemon/internal/server"
	"github.com/stretchr/testify/require"
)

var cfg = c.Config{
	StatsParams: c.StatsParamsConf{
		OS:          c.OSLinux,
		M:           2,
		N:           2,
		CPU:         true,
		DisksUsage:  true,
		DisksIoStat: true,
		NetStat:     true,
	},
}

var (
	host = "localhost"
	port = 50051
)

func TestServer(t *testing.T) {
	grpcServer := server.NewServer(cfg, host, port)
	ctx := context.Background()
	var err error
	go func() {
		err = grpcServer.Run(ctx)
	}()
	<-time.After(time.Second * 3)
	require.NoError(t, err)
	grpcServer.Stop()
}
