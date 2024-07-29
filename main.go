package main

import (
	"fmt"
	"regexp"
	"strings"
	// "strings"
	"log/slog"
	"strconv"
)

var s = `top - 16:45:22 up  4:29,  0 user,  load average: 0.01, 0.00, 0.00
Tasks:   3 total,   1 running,   2 sleeping,   0 stopped,   0 zombie
%Cpu(s):  0.0 us,  0.1 sy,  0.0 ni,100.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
MiB Mem :  32026.0 total,  24553.4 free,   1603.0 used,   6325.4 buff/cache
MiB Swap:   8192.0 total,   8191.7 free,      0.3 used.  30423.0 avail Mem

  PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND
    1 root      20   0    4180   3220   2944 S   0.0   0.0   0:00.00 bash
    7 root      20   0    4180   3456   2936 S   0.0   0.0   0:00.01 bash
 3186 root      20   0    8420   4644   2780 R   0.0   0.0   0:00.00 top`

var (
	reLoadAverage *regexp.Regexp
	cpuStats      *regexp.Regexp
)

func init() {
	reLoadAverage = regexp.MustCompile(`load average:\s([\d\.]+)`)
	cpuStats = regexp.MustCompile(`(\d+\.\d+ (us|sy|id))`)
}

func main() {
	fmt.Println(s)
	lines := strings.Split(s, "\n")
	for _, l := range lines {
		fmt.Println(l)
	}

	s1 := parseLoadAverage(lines[0])
	fmt.Println(s1)
	a, b, c := parseCPUStats(lines[2])
	fmt.Println(a, b, c)
}

func parseLoadAverage(s string) float64 {
	match := reLoadAverage.FindStringSubmatch(s)
	if len(match) < 2 {
		slog.Error("parseLoadAverage match error")
	}

	res, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		slog.Error(err.Error())
	}
	return res
}

func parseCPUStats(s string) (float64, float64, float64) {
	res := [3]float64{}
	match := cpuStats.FindAllStringSubmatch(s, -1)
	if len(match) < 3 {
		slog.Error("parseCPUStats match error")
	}
	for i, sub := range match {
		fmt.Println(sub[1])
		spl := strings.Split(sub[1], " ")
		fl, err := strconv.ParseFloat(spl[0], 64)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		res[i] = fl

	}
	return res[0], res[1], res[2]
}
