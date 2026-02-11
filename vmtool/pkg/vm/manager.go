package vm

import (
	"context"
	"fmt"
	"sync"

	"github.com/utmapp/vmtool/pkg/config"
	"github.com/utmapp/vmtool/pkg/qemu"
)

type Manager struct {
	store    *Store
	running  map[string]*qemu.Runner
	mu       sync.Mutex
}

func NewManager(store *Store) *Manager {
	return &Manager{
		store:   store,
		running: make(map[string]*qemu.Runner),
	}
}

func (m *Manager) StartVM(ctx context.Context, name string) error {
	m.mu.Lock()
	if _, ok := m.running[name]; ok {
		m.mu.Unlock()
		return fmt.Errorf("VM %s is already running", name)
	}
	// Reserve the slot while holding the lock to prevent concurrent starts.
	m.running[name] = nil
	m.mu.Unlock()

	cfg, ok := m.store.GetVM(name)
	if !ok {
		// Clean up reservation if the VM is not found.
		m.mu.Lock()
		delete(m.running, name)
		m.mu.Unlock()
		return fmt.Errorf("VM %s not found", name)
	}

	runner := qemu.NewRunner(cfg)
	if err := runner.Start(ctx); err != nil {
		// Clean up reservation on start failure.
		m.mu.Lock()
		delete(m.running, name)
		m.mu.Unlock()
		return err
	}

	m.mu.Lock()
	m.running[name] = runner
	m.mu.Unlock()

	go func() {
		runner.Wait()
		m.mu.Lock()
		delete(m.running, name)
		m.mu.Unlock()
	}()

	return nil
}

func (m *Manager) StopVM(name string) error {
	m.mu.Lock()
	runner, ok := m.running[name]
	m.mu.Unlock()

	if !ok {
		return fmt.Errorf("VM %s is not running", name)
	}

	return runner.Stop()
}

func (m *Manager) PauseVM(name string) error {
	m.mu.Lock()
	runner, ok := m.running[name]
	m.mu.Unlock()

	if !ok {
		return fmt.Errorf("VM %s is not running", name)
	}

	return runner.Pause()
}

func (m *Manager) ResumeVM(name string) error {
	m.mu.Lock()
	runner, ok := m.running[name]
	m.mu.Unlock()

	if !ok {
		return fmt.Errorf("VM %s is not running", name)
	}

	return runner.Resume()
}

func (m *Manager) CreateSnapshot(vmName string, snapName string) error {
	m.mu.Lock()
	runner, ok := m.running[vmName]
	m.mu.Unlock()

	if !ok {
		return fmt.Errorf("VM %s is not running", vmName)
	}

	return runner.CreateSnapshot(snapName)
}

func (m *Manager) RestoreSnapshot(vmName string, snapName string) error {
	m.mu.Lock()
	runner, ok := m.running[vmName]
	m.mu.Unlock()

	if !ok {
		return fmt.Errorf("VM %s is not running", vmName)
	}

	return runner.RestoreSnapshot(snapName)
}

func (m *Manager) DeleteSnapshot(vmName string, snapName string) error {
	m.mu.Lock()
	runner, ok := m.running[vmName]
	m.mu.Unlock()

	if !ok {
		return fmt.Errorf("VM %s is not running", vmName)
	}

	return runner.DeleteSnapshot(snapName)
}

func (m *Manager) GetStatus(name string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.running[name]; ok {
		return "running"
	}
	return "stopped"
}

func (m *Manager) GetVNCPort(name string) int {
	m.mu.Lock()
	runner, ok := m.running[name]
	m.mu.Unlock()

	if !ok {
		return 0
	}

	return runner.GetVNCPort()
}

func (m *Manager) ListVMs() []*config.VMConfig {
	return m.store.ListVMs()
}
