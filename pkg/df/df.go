package df

import (
	"os/exec"
)

func DiskFreeStats() (string, error) {
	cmd := exec.Command("df", "-k")
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}

func DiskFreeInodeStats() (string, error) {
	cmd := exec.Command("df", "-i")
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}
