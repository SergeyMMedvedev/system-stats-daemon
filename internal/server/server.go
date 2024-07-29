package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/SergeyMMedvedev/system-stats-daemon/internal/collector"
	"github.com/SergeyMMedvedev/system-stats-daemon/internal/config"
	"github.com/SergeyMMedvedev/system-stats-daemon/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedSystemStatsServiceServer
	srv  *grpc.Server
	cfg  config.Config
	host string
	port int
}

func (s *Server) StreamSystemStats(
	_ *pb.SystemStatsRequest,
	stream pb.SystemStatsService_StreamSystemStatsServer,
) error {
	statsCollector := collector.NewCollector(s.cfg)
	statsCollector.CollectStats(context.Background())
	for {
		resp := &pb.SystemStatsResponse{}
		if s.cfg.StatsParams.CPU {
			resp.Cpu = &pb.CPU{}
			cpuStats := <-statsCollector.CpuChan
			resp.LoadAverage = cpuStats[0]
			resp.Cpu.UserMode = cpuStats[1]
			resp.Cpu.SystemMode = cpuStats[2]
			resp.Cpu.Idle = cpuStats[3]
		}
		if s.cfg.StatsParams.DisksUsage {
			resp.Disks = make([]*pb.Disk, 0)
			disksStats := <-statsCollector.DisksFreeChan
			for _, diskStats := range disksStats {
				resp.Disks = append(resp.Disks, &pb.Disk{
					Mounted: diskStats.MountedOn,
					Use:     diskStats.Use,
					Iuse:    diskStats.IUse,
				})
			}
		}
		if s.cfg.StatsParams.OS == config.OSLinux {
			if s.cfg.StatsParams.DisksIoStat {
				resp.Iostat = make([]*pb.IoStat, 0)
				disksIoStats := <-statsCollector.DisksIoStatChan
				for _, diskIoStats := range disksIoStats {
					resp.Iostat = append(resp.Iostat, &pb.IoStat{
						DiskDevice: diskIoStats.DiskDevice,
						Tps:        diskIoStats.Tps,
						KbRead:     diskIoStats.KBReadPS,
						KbWrtn:     diskIoStats.KBWrtnPS,
					})
				}
			}
			if s.cfg.StatsParams.NetStat {
				resp.Netstat = make([]*pb.NetStat, 0)
				netStats := <-statsCollector.NetStatChan
				for _, netStat := range netStats {
					resp.Netstat = append(resp.Netstat, &pb.NetStat{
						Program:  netStat.ProgramName,
						Pid:      int32(netStat.Pid),
						User:     netStat.User,
						Protocol: netStat.Protocol,
						Port:     int32(netStat.Port),
						State:    netStat.State,
					})
				}

				tcpStats := <-statsCollector.TCPStatChan
				resp.Tcpstat = &pb.TCPStat{
					All:      int32(tcpStats.All),
					Estab:    int32(tcpStats.Estab),
					Closed:   int32(tcpStats.Closed),
					Orphaned: int32(tcpStats.Orphaned),
					Timewait: int32(tcpStats.Timewait),
				}
			}
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
}

func (s *Server) Run(_ context.Context) error {
	slog.Info("Starting server")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		return err
	}

	s.srv = grpc.NewServer()
	pb.RegisterSystemStatsServiceServer(s.srv, s)
	reflection.Register(s.srv)
	slog.Info(fmt.Sprintf("Server started on %v:%v", s.host, s.port))
	return s.srv.Serve(lis)
}

func (s *Server) Stop() {
	slog.Info("Stopping server")
	s.srv.GracefulStop()
}

func NewServer(cfg config.Config, host string, port int) *Server {
	return &Server{
		UnimplementedSystemStatsServiceServer: pb.UnimplementedSystemStatsServiceServer{},
		cfg:                                   cfg,
		port:                                  port,
		host:                                  host,
	}
}
