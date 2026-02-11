package qemu

import (
	"strings"
	"testing"

	"github.com/utmapp/vmtool/pkg/config"
)

func TestBuildArgs(t *testing.T) {
	cfg := &config.VMConfig{
		Name: "test-vm",
		UUID: "1234",
		System: config.SystemConfig{
			Architecture: "x86_64",
			Memory:       1024,
			CPUs:         1,
		},
		Display: config.DisplayConfig{
			Enabled: true,
			VNCAddr: ":1",
		},
	}
	builder := NewBuilder(cfg)
	args := builder.BuildArgs()

	joined := strings.Join(args, " ")
	expected := []string{"-name test-vm", "-uuid 1234", "-m 1024", "-smp cpus=1", "-vnc :1"}

	for _, exp := range expected {
		if !strings.Contains(joined, exp) {
			t.Errorf("expected argument %s not found in %s", exp, joined)
		}
	}
}
