package netstat

import (
	"os/exec"
)

func NetStat() (string, error) {
	cmd := exec.Command("sudo", "netstat", "-ltupe")
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}
