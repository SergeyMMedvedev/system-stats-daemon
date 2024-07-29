package top

import (
	"os/exec"
)

func Top() (string, error) {
	cmd := exec.Command("top", "-b", "-n1")
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}
