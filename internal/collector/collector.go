package collector

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/SergeyMMedvedev/system-stats-daemon/internal/collector/cpu"
	"github.com/SergeyMMedvedev/system-stats-daemon/internal/collector/df"
	"github.com/SergeyMMedvedev/system-stats-daemon/internal/collector/iostat"
	"github.com/SergeyMMedvedev/system-stats-daemon/internal/collector/netstat"
	"github.com/SergeyMMedvedev/system-stats-daemon/internal/collector/tcpstat"
	"github.com/SergeyMMedvedev/system-stats-daemon/internal/config"
	"github.com/SergeyMMedvedev/system-stats-daemon/internal/ringbuffer"
)

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

type DiskIoStatBuf struct {
	DiskDevice string
	Tps        ringbuffer.RingBuffer
	KBReadPS   ringbuffer.RingBuffer
	KBWrtnPS   ringbuffer.RingBuffer
}

type NetStat struct {
	Pid         int
	ProgramName string
	User        string
	Protocol    string
	Port        int
	State       string
}

type TCPStatsBuf struct {
	All      ringbuffer.RingBuffer
	Estab    ringbuffer.RingBuffer
	Closed   ringbuffer.RingBuffer
	Orphaned ringbuffer.RingBuffer
	Timewait ringbuffer.RingBuffer
}

type TCPStats struct {
	All      float64
	Estab    float64
	Closed   float64
	Orphaned float64
	Timewait float64
}

type Collector struct {
	cfg config.Config

	CPUMu              *sync.Mutex
	LoadAverageStats   ringbuffer.RingBuffer
	CPUUserModeStats   ringbuffer.RingBuffer
	CPUSystemModeStats ringbuffer.RingBuffer
	CPUIdleStats       ringbuffer.RingBuffer

	CPUChan chan [4]float64

	DisksFreeStatsMu    *sync.Mutex
	DisksFreeStats      []ringbuffer.RingBuffer
	DisksFreeInodeStats []ringbuffer.RingBuffer
	DisksFreeChan       chan []DiskFree

	DisksIoStatMu   *sync.Mutex
	DisksIoStatBuf  []DiskIoStatBuf
	DisksIoStatChan chan []DiskIoStat

	NetStatChan chan []NetStat

	TCPStatsMu  *sync.Mutex
	TCPStatsBuf TCPStatsBuf
	TCPStatChan chan TCPStats
}

func NewCollector(cfg config.Config) *Collector {
	m := cfg.StatsParams.M
	return &Collector{
		cfg:                cfg,
		LoadAverageStats:   *ringbuffer.NewRingBuffer(m),
		CPUUserModeStats:   *ringbuffer.NewRingBuffer(m),
		CPUSystemModeStats: *ringbuffer.NewRingBuffer(m),
		CPUIdleStats:       *ringbuffer.NewRingBuffer(m),
		TCPStatsBuf: TCPStatsBuf{
			All:      *ringbuffer.NewRingBuffer(m),
			Estab:    *ringbuffer.NewRingBuffer(m),
			Closed:   *ringbuffer.NewRingBuffer(m),
			Orphaned: *ringbuffer.NewRingBuffer(m),
			Timewait: *ringbuffer.NewRingBuffer(m),
		},
		DisksIoStatMu:    &sync.Mutex{},
		TCPStatsMu:       &sync.Mutex{},
		DisksFreeStatsMu: &sync.Mutex{},
		CPUMu:            &sync.Mutex{},
	}
}

func (c *Collector) collectCPUstats(ctx context.Context, cfg config.Config) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("Stop collect cpu stats")
			return
		default:
			loadAverage, userMode, systemMode, idle, err := cpu.CollectCPUStats(cfg.StatsParams.OS)
			if err != nil {
				slog.Error("collect CPU stats err:" + err.Error())
				continue
			}
			c.CPUMu.Lock()
			c.LoadAverageStats.Enqueue(loadAverage)
			c.CPUUserModeStats.Enqueue(userMode)
			c.CPUSystemModeStats.Enqueue(systemMode)
			c.CPUIdleStats.Enqueue(idle)
			c.CPUMu.Unlock()
			time.Sleep(time.Second)
		}
	}
}

func (c *Collector) sendCPUstats(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("Stop send cpu stats")
			return
		default:
			c.CPUMu.Lock()
			if c.CPUUserModeStats.IsFull() && c.CPUSystemModeStats.IsFull() && c.CPUIdleStats.IsFull() {
				slog.Info("send CPUstats")
				cpu := [4]float64{
					c.LoadAverageStats.Average(),
					c.CPUUserModeStats.Average(),
					c.CPUSystemModeStats.Average(),
					c.CPUIdleStats.Average(),
				}
				c.CPUChan <- cpu
			}
			c.CPUMu.Unlock()
			time.Sleep(time.Second * time.Duration(c.cfg.StatsParams.N))
		}
	}
}

func (c *Collector) collectDiskFreeStats(ctx context.Context) {
	disks, err := df.CollectDiskFreeStats()
	if err != nil {
		slog.Error(err.Error())
	}
	disksNum := len(disks)
	for i := 0; i < disksNum; i++ {
		c.DisksFreeStatsMu.Lock()
		c.DisksFreeStats = append(c.DisksFreeStats, *ringbuffer.NewRingBuffer(c.cfg.StatsParams.M))
		c.DisksFreeInodeStats = append(c.DisksFreeInodeStats, *ringbuffer.NewRingBuffer(c.cfg.StatsParams.M))
		c.DisksFreeStatsMu.Unlock()
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stop collect disks stats")
			return
		default:
			disks, err := df.CollectDiskFreeStats()
			if err != nil {
				slog.Error(err.Error())
			}
			disksInode, err := df.CollectDiskFreeInodeStats()
			if err != nil {
				slog.Error(err.Error())
			}
			for i := 0; i < disksNum; i++ {
				c.DisksFreeStatsMu.Lock()
				c.DisksFreeStats[i].Enqueue(disks[i].Use)
				c.DisksFreeInodeStats[i].Enqueue(disksInode[i].Use)
				c.DisksFreeStatsMu.Unlock()
			}
			time.Sleep(time.Second)
		}
	}
}

func (c *Collector) sendDiskFreeStats(ctx context.Context) {
	disks, err := df.CollectDiskFreeStats()
	if err != nil {
		slog.Error(err.Error())
	}
	disksNum := len(disks)
	for {
		select {
		case <-ctx.Done():
			slog.Info("Stop collect disks stats")
			return
		default:
			disksFree := make([]DiskFree, 0)
			c.DisksFreeStatsMu.Lock()
			if c.DisksFreeStats != nil && c.DisksFreeInodeStats != nil {
				for i := 0; i < disksNum; i++ {
					if c.DisksFreeStats[i].IsFull() && c.DisksFreeInodeStats[i].IsFull() {
						disksFree = append(disksFree, DiskFree{
							MountedOn: disks[i].MountedOn,
							Use:       fmt.Sprintf("%f", c.DisksFreeStats[i].Average()) + "%",
							IUse:      fmt.Sprintf("%f", c.DisksFreeInodeStats[i].Average()) + "%",
						})
					}
				}
				slog.Info("send disksFree stats")
				c.DisksFreeChan <- disksFree
			}
			c.DisksFreeStatsMu.Unlock()
			time.Sleep(time.Second * time.Duration(c.cfg.StatsParams.N))
		}
	}
}

func (c *Collector) collectDiskIoStat(ctx context.Context) {
	disks, err := iostat.CollectIoStat()
	if err != nil {
		slog.Error(err.Error())
	}
	disksNum := len(disks)
	for i := 0; i < disksNum; i++ {
		d := DiskIoStatBuf{
			Tps:      *ringbuffer.NewRingBuffer(c.cfg.StatsParams.M),
			KBReadPS: *ringbuffer.NewRingBuffer(c.cfg.StatsParams.M),
			KBWrtnPS: *ringbuffer.NewRingBuffer(c.cfg.StatsParams.M),
		}
		c.DisksIoStatMu.Lock()
		c.DisksIoStatBuf = append(c.DisksIoStatBuf, d)
		c.DisksIoStatMu.Unlock()
	}
	for {
		select {
		case <-ctx.Done():
			slog.Info("Stop collect iostats")
			return
		default:
			disks, err := iostat.CollectIoStat()
			if err != nil {
				slog.Error(err.Error())
			}
			for i := 0; i < disksNum; i++ {
				c.DisksIoStatMu.Lock()
				c.DisksIoStatBuf[i].DiskDevice = disks[i].DiskDevice
				c.DisksIoStatBuf[i].Tps.Enqueue(disks[i].Tps)
				c.DisksIoStatBuf[i].KBReadPS.Enqueue(disks[i].KBReadPS)
				c.DisksIoStatBuf[i].KBWrtnPS.Enqueue(disks[i].KBWrtnPS)
				c.DisksIoStatMu.Unlock()
			}
			time.Sleep(time.Second)
		}
	}
}

func (c *Collector) sendDiskIoStat(ctx context.Context) {
	disks, err := iostat.CollectIoStat()
	if err != nil {
		slog.Error(err.Error())
	}
	disksNum := len(disks)
	for {
		select {
		case <-ctx.Done():
			slog.Info("Stop collect iostats")
			return
		default:
			disksIoStat := make([]DiskIoStat, 0)
			c.DisksIoStatMu.Lock()
			if c.DisksIoStatBuf != nil {
				for i := 0; i < disksNum; i++ {
					if len(c.DisksIoStatBuf) == disksNum &&
						c.DisksIoStatBuf[i].Tps.IsFull() &&
						c.DisksIoStatBuf[i].KBReadPS.IsFull() &&
						c.DisksIoStatBuf[i].KBWrtnPS.IsFull() {
						disksIoStat = append(disksIoStat, DiskIoStat{
							DiskDevice: c.DisksIoStatBuf[i].DiskDevice,
							Tps:        c.DisksIoStatBuf[i].Tps.Average(),
							KBReadPS:   c.DisksIoStatBuf[i].KBReadPS.Average(),
							KBWrtnPS:   c.DisksIoStatBuf[i].KBWrtnPS.Average(),
						})
					}
				}
				slog.Info("send disks Io stats")
				c.DisksIoStatChan <- disksIoStat
			}
			c.DisksIoStatMu.Unlock()
			time.Sleep(time.Second * time.Duration(c.cfg.StatsParams.N))
		}
	}
}

func (c *Collector) collectAndSendNetStat(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("Stop collect netstat")
			return
		default:
			netStats, err := netstat.CollectNetstat()
			if err != nil {
				slog.Error(err.Error())
			}
			netStatsResult := make([]NetStat, 0)
			for _, netstat := range netStats {
				n := NetStat{
					Pid:         netstat.PID,
					ProgramName: netstat.Program,
					User:        netstat.User,
					Protocol:    netstat.Proto,
					Port:        netstat.LocalPort,
					State:       netstat.State,
				}
				netStatsResult = append(netStatsResult, n)
			}
			slog.Info("send Net Stat")
			c.NetStatChan <- netStatsResult
			time.Sleep(time.Second)
		}
	}
}

func (c *Collector) collectTCPStat(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("Stop collect netstat")
			return
		default:
			tcpstats, err := tcpstat.CollectTCPStats()
			if err != nil {
				slog.Error(err.Error())
			}
			c.TCPStatsMu.Lock()
			c.TCPStatsBuf.All.Enqueue(float64(tcpstats.All))
			c.TCPStatsBuf.Estab.Enqueue(float64(tcpstats.Estab))
			c.TCPStatsBuf.Closed.Enqueue(float64(tcpstats.Closed))
			c.TCPStatsBuf.Orphaned.Enqueue(float64(tcpstats.Orphaned))
			c.TCPStatsBuf.Timewait.Enqueue(float64(tcpstats.Timewait))
			c.TCPStatsMu.Unlock()
			time.Sleep(time.Second)
		}
	}
}

func (c *Collector) sendTCPStat(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("Stop collect netstat")
			return
		default:
			c.TCPStatsMu.Lock()
			if c.TCPStatsBuf.All.IsFull() &&
				c.TCPStatsBuf.Estab.IsFull() &&
				c.TCPStatsBuf.Closed.IsFull() &&
				c.TCPStatsBuf.Orphaned.IsFull() &&
				c.TCPStatsBuf.Timewait.IsFull() {
				t := TCPStats{
					All:      c.TCPStatsBuf.All.Average(),
					Estab:    c.TCPStatsBuf.Estab.Average(),
					Closed:   c.TCPStatsBuf.Closed.Average(),
					Orphaned: c.TCPStatsBuf.Orphaned.Average(),
					Timewait: c.TCPStatsBuf.Timewait.Average(),
				}
				slog.Info("send TCP Stat")
				c.TCPStatChan <- t
			}
			c.TCPStatsMu.Unlock()
			time.Sleep(time.Second * time.Duration(c.cfg.StatsParams.N))
		}
	}
}

func (c *Collector) CollectStats(ctx context.Context) {
	if c.cfg.StatsParams.CPU {
		c.CPUChan = make(chan [4]float64)
		go c.collectCPUstats(ctx, c.cfg)
		go c.sendCPUstats(ctx)
	}
	if c.cfg.StatsParams.DisksUsage {
		c.DisksFreeChan = make(chan []DiskFree)
		go c.collectDiskFreeStats(ctx)
		go c.sendDiskFreeStats(ctx)
	}
	if c.cfg.StatsParams.OS == config.OSLinux {
		if c.cfg.StatsParams.DisksIoStat {
			c.DisksIoStatChan = make(chan []DiskIoStat)
			go c.collectDiskIoStat(ctx)
			go c.sendDiskIoStat(ctx)
		}
		if c.cfg.StatsParams.NetStat {
			c.NetStatChan = make(chan []NetStat)
			go c.collectAndSendNetStat(ctx)

			c.TCPStatChan = make(chan TCPStats)
			go c.collectTCPStat(ctx)
			go c.sendTCPStat(ctx)
		}
	}
}
