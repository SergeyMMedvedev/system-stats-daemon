package dpkg

import (
	"os/exec"
	"strings"
)

func CheckInstall(packageName string) (bool, error) {
	cmd := exec.Command("dpkg", "-s", packageName)
	stdout, err := cmd.Output()
	if err != nil {
		return false, err
	}
	result := strings.Contains(string(stdout), "Status: install ok installed")
	return result, nil
}
