package cpu

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"github.com/SergeyMMedvedev/system-stats-daemon/internal/config"
	"github.com/SergeyMMedvedev/system-stats-daemon/pkg/top"
	"github.com/SergeyMMedvedev/system-stats-daemon/pkg/wmic"
)

var (
	reLoadAverage *regexp.Regexp
	cpuStats      *regexp.Regexp
)

func init() {
	reLoadAverage = regexp.MustCompile(`load average:\s([\d\.]+)`)
	cpuStats = regexp.MustCompile(`(\d+\.\d+ (us|sy|id),)`)
}

func parseLoadAverage(s string) (float64, error) {
	match := reLoadAverage.FindStringSubmatch(s)
	if len(match) < 2 {
		return 0, fmt.Errorf("parseLoadAverage match error")
	}
	res, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0, err
	}
	return res, nil
}

type CPUStats struct {
	UserMode   float64
	SystemMode float64
	Idle       float64
}

func parseCPUStats(s string) (CPUStats, error) {
	cpuStatsArr := []CPUStats{}
	match := cpuStats.FindAllStringSubmatch(s, -1)
	if len(match) < 3 {
		return CPUStats{}, fmt.Errorf("parseCPUStats match error")
	}
	for i := 0; i < len(match); i += 3 {
		cpuStats := CPUStats{}
		for j := i; j < i+3; j++ {
			sub := match[j]
			spl := strings.Split(sub[1], " ")
			cpuMode := spl[1]
			value, err := strconv.ParseFloat(spl[0], 64)
			if err != nil {
				slog.Error(err.Error())
				continue
			}
			switch {
			case strings.HasPrefix(cpuMode, "us"):
				cpuStats.UserMode = value
			case strings.HasPrefix(cpuMode, "sy"):
				cpuStats.SystemMode = value
			case strings.HasPrefix(cpuMode, "id"):
				cpuStats.Idle = value
			}
		}
		cpuStatsArr = append(cpuStatsArr, cpuStats)
	}
	cpuStatsMinIdle := cpuStatsArr[0]
	for _, cpuStats := range cpuStatsArr {
		if cpuStats.Idle < cpuStatsMinIdle.Idle {
			cpuStatsMinIdle = cpuStats
		}
	}
	return cpuStatsMinIdle, nil
}

func collectLinuxCPUStats() (loadAvg float64, cpuStats CPUStats, err error) {
	cpuStatsStr, err := top.Top()
	if err != nil {
		slog.Error(err.Error())
	}
	lines := strings.Split(cpuStatsStr, "\n")

	loadAvg, err = parseLoadAverage(lines[0])
	if err != nil {
		return 0, cpuStats, err
	}
	cpuStats, err = parseCPUStats(cpuStatsStr)
	if err != nil {
		return 0, cpuStats, err
	}
	return loadAvg, cpuStats, nil
}

func collectWindowsCPUStats() (loadAvg float64, cpuStats CPUStats, err error) {
	p, err := wmic.CPUGetLoadPercentage()
	if err != nil {
		return 0, cpuStats, err
	}
	cpuStats.UserMode = float64(p)
	cpuStats.SystemMode = float64(p)
	cpuStats.Idle = float64(100 - p)
	return 0, cpuStats, nil
}

func CollectCPUStats(os config.OS) (loadAvg float64, cpuStats CPUStats, err error) {
	if os == config.OSWindows {
		return collectWindowsCPUStats()
	}
	return collectLinuxCPUStats()
}
