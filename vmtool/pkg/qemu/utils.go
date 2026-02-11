package qemu

import (
	"fmt"
	"os/exec"
)

func CreateDisk(path string, size string, format string) error {
	if format == "" {
		format = "qcow2"
	}
	cmd := exec.Command("qemu-img", "create", "-f", format, path, size)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create disk: %v, output: %s", err, string(output))
	}
	return nil
}
