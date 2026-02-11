package vm

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/utmapp/vmtool/pkg/config"
	"gopkg.in/yaml.v3"
)

type Store struct {
	baseDir string
	vms     map[string]*config.VMConfig
	mu      sync.RWMutex
}

func NewStore(baseDir string) (*Store, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	s := &Store{
		baseDir: baseDir,
		vms:     make(map[string]*config.VMConfig),
	}
	if err := s.LoadAll(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) LoadAll() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	files, err := os.ReadDir(s.baseDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" {
			cfg, err := s.loadVM(filepath.Join(s.baseDir, f.Name()))
			if err != nil {
				fmt.Printf("Warning: failed to load VM config %s: %v\n", f.Name(), err)
				continue
			}
			s.vms[cfg.Name] = cfg
		}
	}
	return nil
}

func (s *Store) loadVM(path string) (*config.VMConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg config.VMConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (s *Store) SaveVM(cfg *config.VMConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	path := filepath.Join(s.baseDir, cfg.Name+".yaml")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	s.vms[cfg.Name] = cfg
	return nil
}

func (s *Store) GetVM(name string) (*config.VMConfig, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cfg, ok := s.vms[name]
	return cfg, ok
}

func (s *Store) ListVMs() []*config.VMConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var list []*config.VMConfig
	for _, v := range s.vms {
		list = append(list, v)
	}
	return list
}

func (s *Store) DeleteVM(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.baseDir, name+".yaml")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	delete(s.vms, name)
	return nil
}
