package qemu

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/utmapp/vmtool/pkg/config"
)

type Runner struct {
	config  *config.VMConfig
	cmd     *exec.Cmd
	cancel  context.CancelFunc
}

func NewRunner(cfg *config.VMConfig) *Runner {
	return &Runner{config: cfg}
}

func (r *Runner) Start(ctx context.Context) error {
	builder := NewBuilder(r.config)
	args := builder.BuildArgs()

	// Find qemu binary
	qemuBin := r.findQemuBinary()

	// Add QMP support
	qmpSocket := r.getQMPSocketPath()
	args = append(args, "-qmp", "unix:"+qmpSocket+",server,nowait")

	r.cmd = exec.CommandContext(ctx, qemuBin, args...)
	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stderr

	return r.cmd.Start()
}

func (r *Runner) Stop() error {
	client := NewQMPClient(r.getQMPSocketPath())
	if err := client.PowerDown(); err == nil {
		return nil
	}
	// Fallback to kill if QMP fails or isn't responsive
	if r.cmd != nil && r.cmd.Process != nil {
		return r.cmd.Process.Signal(syscall.SIGTERM)
	}
	return nil
}

func (r *Runner) Pause() error {
	return NewQMPClient(r.getQMPSocketPath()).Pause()
}

func (r *Runner) Resume() error {
	return NewQMPClient(r.getQMPSocketPath()).Resume()
}

func (r *Runner) CreateSnapshot(name string) error {
	return NewQMPClient(r.getQMPSocketPath()).SaveSnapshot(name)
}

func (r *Runner) RestoreSnapshot(name string) error {
	return NewQMPClient(r.getQMPSocketPath()).LoadSnapshot(name)
}

func (r *Runner) DeleteSnapshot(name string) error {
	return NewQMPClient(r.getQMPSocketPath()).DeleteSnapshot(name)
}

func (r *Runner) GetVNCPort() int {
	return r.config.Display.VNCPort
}

func (r *Runner) Wait() error {
	if r.cmd != nil {
		return r.cmd.Wait()
	}
	return nil
}

func (r *Runner) findQemuBinary() string {
	arch := r.config.System.Architecture
	if arch == "" {
		arch = "x86_64"
	}
	binName := "qemu-system-" + arch
	if path, err := exec.LookPath(binName); err == nil {
		return path
	}
	return binName // Fallback to raw name
}

func (r *Runner) getQMPSocketPath() string {
	dataDir := config.GetDefaultDataDir()
	return filepath.Join(dataDir, r.config.UUID+".qmp")
}
