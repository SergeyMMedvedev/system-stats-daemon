package apt

import (
	"log/slog"
	"os/exec"
)

func Update() error {
	cmd := exec.Command("apt-get", "update")
	stdout, err := cmd.Output()
	if err != nil {
		return err
	}
	slog.Info(string(stdout))
	return nil
}

func Install(packageName string) error {
	cmd := exec.Command(
		"apt-get", "install", "-y", "--no-install-recommends", packageName,
	)
	stdout, err := cmd.Output()
	if err != nil {
		return err
	}
	slog.Info(string(stdout))
	return nil
}
