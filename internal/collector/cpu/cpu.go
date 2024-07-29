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
	cpuStats = regexp.MustCompile(`(\d+\.\d+ (us|sy|id))`)
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

func parseCPUStats(s string) (float64, float64, float64, error) {
	res := [3]float64{}
	match := cpuStats.FindAllStringSubmatch(s, -1)
	if len(match) < 3 {
		return 0, 0, 0, fmt.Errorf("parseCPUStats match error")
	}
	for i, sub := range match {
		spl := strings.Split(sub[1], " ")
		fl, err := strconv.ParseFloat(spl[0], 64)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		res[i] = fl

	}
	return res[0], res[1], res[2], nil
}

func collectLinuxCPUStats() (float64, float64, float64, float64, error) {
	cpuStats, err := top.Top()
	if err != nil {
		slog.Error(err.Error())
	}
	lines := strings.Split(cpuStats, "\n")

	loadAvg, err := parseLoadAverage(lines[0])
	if err != nil {
		return 0, 0, 0, 0, err
	}
	us, sy, id, err := parseCPUStats(lines[2])
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return loadAvg, us, sy, id, nil
}

func collectWindowsCPUStats() (float64, float64, float64, float64, error) {
	p, err := wmic.CPUGetLoadPercentage()
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return 0, float64(p), float64(p), float64(100 - p), nil
}

func CollectCPUStats(os config.OS) (float64, float64, float64, float64, error) {
	if os == config.OSWindows {
		return collectWindowsCPUStats()
	}
	return collectLinuxCPUStats()
}
