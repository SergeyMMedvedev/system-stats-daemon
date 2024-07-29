package wmic

import (
	"os/exec"
	"strconv"
	"strings"
)

func CPUGetLoadPercentage() (int, error) {
	var percentage int
	cmd := exec.Command("wmic", "cpu", "get", "loadpercentage")
	stdout, err := cmd.Output()
	if err != nil {
		return percentage, err
	}
	out := strings.Trim(string(stdout), " \n\t\r")
	arr := strings.Split(out, "\n")
	if len(arr) != 2 {
		return percentage, nil
	}
	percentage, err = strconv.Atoi(arr[1])
	if err != nil {
		return percentage, err
	}
	return percentage, nil
}
