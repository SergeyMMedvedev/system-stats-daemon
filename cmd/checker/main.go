package main

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
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
	fmt.Println("len(match)", len(match))
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
			case cpuMode == "us,":
				cpuStats.UserMode = value
			case cpuMode == "sy,":
				cpuStats.SystemMode = value
			case cpuMode == "id,":
				cpuStats.Idle = value
			}
		}
		cpuStatsArr = append(cpuStatsArr, cpuStats)
	}
	fmt.Println("cpuStatsArr", cpuStatsArr)
	cpuStatsMinIdle := cpuStatsArr[0]
	for _, cpuStats := range cpuStatsArr {
		if cpuStats.Idle < cpuStatsMinIdle.Idle {
			cpuStatsMinIdle = cpuStats
		}
	}
	return cpuStatsMinIdle, nil
}

func main() {
	s := `%Cpu(s):  0.0 us,  0.0 sy,  0.0 ni,100.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st`

	r, err := parseLoadAverage(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)

	cpuStats, err := parseCPUStats(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cpuStats)
}
