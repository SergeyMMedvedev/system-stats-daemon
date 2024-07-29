package netstat

import (
	"bufio"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/SergeyMMedvedev/system-stats-daemon/pkg/netstat"
)

type NetstatEntry struct {
	Proto          string
	RecvQ          int
	SendQ          int
	LocalPort      int
	ForeignAddress string
	State          string
	User           string
	Inode          int
	PID            int
	Program        string
}

func parseNetStat(netstat string) []NetstatEntry {
	fmt.Println(netstat)
	var entries []NetstatEntry

	scanner := bufio.NewScanner(strings.NewReader(netstat))

	if scanner.Scan() {
		// The first line is the header, ignore it
	}
	if scanner.Scan() {
		// The second line is the header, ignore it
	}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		fields := strings.Fields(line)

		if len(fields) < 8 {
			slog.Error("Unexpected format: " + line)
			continue
		}

		proto := fields[0]
		recvQ, err := strconv.Atoi(fields[1])
		if err != nil {
			slog.Error("failed to convert recvQ: " + err.Error())
		}
		sendQ, err := strconv.Atoi(fields[2])
		if err != nil {
			slog.Error("failed to convert sendQ:" + err.Error())
		}
		localAddress := fields[3]
		foreignAddress := fields[4]
		state := fields[5]
		user := fields[6]
		inode, _ := strconv.Atoi(fields[7])
		pidProgram := strings.Join(fields[8:], " ")

		splitLocalAddress := strings.Split(localAddress, ":")
		LocalPort, err := strconv.Atoi(splitLocalAddress[len(splitLocalAddress)-1])
		if err != nil {
			slog.Error("failed to convert LocalPort: " + err.Error())
		}

		splitPidProgram := strings.Split(pidProgram, "/")
		if len(splitPidProgram) != 2 {
			slog.Error(
				"check splitPidProgram. split len: " + err.Error(),
			)
		}
		pid, err := strconv.Atoi(splitPidProgram[0])
		if err != nil {
			slog.Error("failed to convert pid:" + err.Error())
		}
		program := splitPidProgram[1]

		entry := NetstatEntry{
			Proto:          proto,
			RecvQ:          recvQ,
			SendQ:          sendQ,
			LocalPort:      LocalPort,
			ForeignAddress: foreignAddress,
			State:          state,
			User:           user,
			Inode:          inode,
			PID:            pid,
			Program:        program,
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		slog.Error("error reading file: " + err.Error())
	}
	return entries
}

func CollectNetstat() ([]NetstatEntry, error) {
	netstat, err := netstat.NetStat()
	if err != nil {
		return nil, fmt.Errorf("fail to get netstat: %w", err)
	}
	netstatEntries := parseNetStat(netstat)
	return netstatEntries, nil
}
