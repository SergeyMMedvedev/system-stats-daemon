package tcpstat

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SergeyMMedvedev/system-stats-daemon/pkg/ss"
)

type TCPStats struct {
	All      int
	Estab    int
	Closed   int
	Orphaned int
	Timewait int
}

func parseTCPStats(input string) (*TCPStats, error) {
	input = strings.TrimPrefix(input, "TCP:")
	input = strings.TrimSpace(input)

	parts := strings.SplitN(input, "(", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("unexpected input format")
	}

	all, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid All value: %w", err)
	}

	params := strings.TrimSuffix(parts[1], ")")
	fields := strings.Split(params, ",")
	if len(fields) != 4 {
		return nil, fmt.Errorf("unexpected number of parameters")
	}

	stats := &TCPStats{All: all}

	for _, field := range fields {
		field = strings.TrimSpace(field)
		keyValue := strings.SplitN(field, " ", 2)
		if len(keyValue) != 2 {
			return nil, fmt.Errorf("invalid field format: %s", field)
		}

		key := strings.TrimSpace(keyValue[0])
		value, err := strconv.Atoi(strings.TrimSpace(keyValue[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid value for %s: %w", key, err)
		}

		switch key {
		case "estab":
			stats.Estab = value
		case "closed":
			stats.Closed = value
		case "orphaned":
			stats.Orphaned = value
		case "timewait":
			stats.Timewait = value
		default:
			return nil, fmt.Errorf("unknown key: %s", key)
		}
	}

	return stats, nil
}

func CollectTCPStats() (*TCPStats, error) {
	tcpstats, err := ss.GrepTCP()
	if err != nil {
		return nil, fmt.Errorf("failed to Grep TCP stats: %w", err)
	}
	tcpstatsStruct, err := parseTCPStats(tcpstats)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TCP stats: %w", err)
	}
	return tcpstatsStruct, nil
}
