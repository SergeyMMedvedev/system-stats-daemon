package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/SergeyMMedvedev/system-stats-daemon/internal/config"
	"github.com/SergeyMMedvedev/system-stats-daemon/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	cl  pb.SystemStatsServiceClient
	cfg config.Config
}

func (c *Client) Connect(host string, port int, cfg config.Config) error {
	slog.Info(fmt.Sprintf("connect to %s:%d", host, port))
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("did not connect: %v", err.Error())
	}
	c.cl = pb.NewSystemStatsServiceClient(conn)
	c.cfg = cfg
	return nil
}

type CPUStats struct {
	UserMode   float64
	SystemMode float64
	Idle       float64
}

type DiskFree struct {
	MountedOn string
	Use       string
	IUse      string
}

type DiskIoStat struct {
	DiskDevice string  `json:"disk_device"`
	Tps        float64 `json:"tps"`
	KBReadPS   float64 `json:"kB_read/s"`
	KBWrtnPS   float64 `json:"kB_wrtn/s"`
}

type NetStat struct {
	Pid         int    `json:"pid"`
	ProgramName string `json:"program"`
	User        string `json:"user"`
	Protocol    string `json:"protocol"`
	Port        int    `json:"port"`
	State       string `json:"state"`
}

type TCPStats struct {
	All      int
	Estab    int
	Closed   int
	Orphaned int
	Timewait int
}

type Stats struct {
	LoadAverage float64
	CPUStats    CPUStats
	DisksFree   []DiskFree
	DisksIoStat []DiskIoStat
	NetStat     []NetStat
	TCPStats    TCPStats
}

func fillRespCPU(stats *Stats, resp *pb.SystemStatsResponse) {
	stats.LoadAverage = resp.LoadAverage
	if resp.Cpu != nil {
		cpuStats := CPUStats{
			UserMode:   resp.Cpu.UserMode,
			SystemMode: resp.Cpu.SystemMode,
			Idle:       resp.Cpu.Idle,
		}
		stats.CPUStats = cpuStats
	}
}

func fillDisksUsage(stats *Stats, resp *pb.SystemStatsResponse) {
	disksFree := []DiskFree{}
	for _, d := range resp.Disks {
		d := DiskFree{
			MountedOn: d.Mounted,
			Use:       d.Use,
			IUse:      d.Iuse,
		}
		disksFree = append(disksFree, d)
	}
	stats.DisksFree = disksFree
}

func fillIoStats(stats *Stats, resp *pb.SystemStatsResponse) {
	disksIoStats := []DiskIoStat{}
	for _, ioStat := range resp.Iostat {
		io := DiskIoStat{
			DiskDevice: ioStat.DiskDevice,
			Tps:        ioStat.Tps,
			KBReadPS:   ioStat.KbRead,
			KBWrtnPS:   ioStat.KbWrtn,
		}
		disksIoStats = append(disksIoStats, io)
	}
	stats.DisksIoStat = disksIoStats
}

func fillNetStat(stats *Stats, resp *pb.SystemStatsResponse) {
	netStats := []NetStat{}
	for _, netStat := range resp.Netstat {
		nt := NetStat{
			ProgramName: netStat.Program,
			Pid:         int(netStat.Pid),
			User:        netStat.User,
			Protocol:    netStat.Protocol,
			Port:        int(netStat.Port),
			State:       netStat.State,
		}
		netStats = append(netStats, nt)
	}
	stats.NetStat = netStats
	if resp.Tcpstat != nil {
		tcpstats := TCPStats{
			All:      int(resp.Tcpstat.All),
			Estab:    int(resp.Tcpstat.Estab),
			Closed:   int(resp.Tcpstat.Closed),
			Orphaned: int(resp.Tcpstat.Orphaned),
			Timewait: int(resp.Tcpstat.Timewait),
		}
		stats.TCPStats = tcpstats
	}
}

func (c *Client) StreamSystemStats(ctx context.Context, req *pb.SystemStatsRequest) error {
	stream, err := c.cl.StreamSystemStats(ctx, req)
	if err != nil {
		return fmt.Errorf("StreamSystemStats err: %v", err.Error())
	}
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			slog.Error(fmt.Sprintf("error receiving response: %v", err))
		}
		stats := &Stats{}
		if c.cfg.StatsParams.CPU {
			fillRespCPU(stats, resp)
		}
		if c.cfg.StatsParams.DisksUsage {
			fillDisksUsage(stats, resp)
		}
		if c.cfg.StatsParams.DisksIoStat {
			fillIoStats(stats, resp)
		}
		if c.cfg.StatsParams.NetStat {
			fillNetStat(stats, resp)
		}
		s, err := json.MarshalIndent(stats, "", "  ")
		if err != nil {
			slog.Error(err.Error())
		}
		fmt.Println(string(s))
		time.Sleep(time.Duration(c.cfg.StatsParams.N))
	}
	return nil
}

func NewClient() *Client {
	return &Client{}
}
