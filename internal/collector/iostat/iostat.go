package iostat

import (
	"encoding/json"
	"fmt"

	"github.com/SergeyMMedvedev/system-stats-daemon/pkg/iostat"
)

type Stat struct {
	Sysstat Sysstat `json:"sysstat"`
}

type Sysstat struct {
	Hosts []Host `json:"hosts"`
}

type Host struct {
	Nodename     string      `json:"nodename"`
	Sysname      string      `json:"sysname"`
	Release      string      `json:"release"`
	Machine      string      `json:"machine"`
	NumberOfCPUs int         `json:"number-of-cpus"`
	Date         string      `json:"date"`
	Statistics   []Statistic `json:"statistics"`
}

type Statistic struct {
	Disk []Disk `json:"disk"`
}

type Disk struct {
	DiskDevice string  `json:"disk_device"`
	Tps        float64 `json:"tps"`
	KBReadPS   float64 `json:"kB_read/s"`
	KBWrtnPS   float64 `json:"kB_wrtn/s"`
	KBRead     int     `json:"kB_read"`
	KBWrtn     int     `json:"kB_wrtn"`
}

func collectSysstat() (*Stat, error) {
	stats, err := iostat.GetStat()
	if err != nil {
		return nil, fmt.Errorf("fail to get sysstat: %w", err)
	}
	var stat Stat
	err = json.Unmarshal(stats, &stat)
	if err != nil {
		return nil, fmt.Errorf("fail to parse sysstat: %w", err)
	}
	return &stat, nil
}

func CollectIoStat() ([]Disk, error) {
	stats, err := collectSysstat()
	if err != nil {
		return nil, fmt.Errorf("fail to get iostat: %w", err)
	}
	hosts := stats.Sysstat.Hosts
	if len(hosts) == 0 {
		return nil, fmt.Errorf("systats hosts length is 0")
	}
	statistics := hosts[0].Statistics
	if len(statistics) == 0 {
		return nil, fmt.Errorf("systats statistics length is 0")
	}
	return statistics[0].Disk, nil
}
