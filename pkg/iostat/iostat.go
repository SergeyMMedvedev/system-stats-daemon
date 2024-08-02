package iostat

import (
	"os/exec"
)

func GetStat() ([]byte, error) {
	cmd := exec.Command("iostat", "-d", "-k", "-o", "JSON")
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return stdout, nil
}
