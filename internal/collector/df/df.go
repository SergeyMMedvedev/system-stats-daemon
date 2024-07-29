package df

import (
	_ "fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"github.com/SergeyMMedvedev/system-stats-daemon/pkg/df"
)

var re *regexp.Regexp

func init() {
	re = regexp.MustCompile(`\s+`)
}

func parseDf(s string) ([]DiskInfo, error) {
	lines := strings.Split(s, "\n")
	result := make([]DiskInfo, 0)
	for _, l := range lines[1:] {
		splitLine := re.Split(l, -1)
		if len(splitLine) >= 2 {
			use := splitLine[len(splitLine)-2]
			useDigit := strings.Trim(use, "%")
			var userInt int
			var err error
			if useDigit == "-" {
				userInt = 0
			} else {
				userInt, err = strconv.Atoi(useDigit)
				if err != nil {
					slog.Error(err.Error())
				}
			}
			userFloat := float64(userInt)
			diskInfo := DiskInfo{
				MountedOn: splitLine[len(splitLine)-1],
				Use:       userFloat,
			}

			result = append(result, diskInfo)
		}
	}
	return result, nil
}

type DiskInfo struct {
	MountedOn string
	Use       float64
}

func CollectDiskFreeStats() ([]DiskInfo, error) {
	stdout, err := df.DiskFreeStats()
	if err != nil {
		return nil, err
	}
	return parseDf(stdout)
}

func CollectDiskFreeInodeStats() ([]DiskInfo, error) {
	stdout, err := df.DiskFreeInodeStats()
	if err != nil {
		return nil, err
	}
	return parseDf(stdout)
}
