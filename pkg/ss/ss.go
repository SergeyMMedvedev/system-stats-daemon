package ss

import (
	"os/exec"
)

func GrepTCP() (string, error) {
	cmd := exec.Command("sh", "-c", "ss -s | grep 'TCP:'")
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}
