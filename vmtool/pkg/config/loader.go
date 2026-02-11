package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadAppConfig() (*AppConfig, error) {
	configDir := GetDefaultConfigDir()
	cfgPath := filepath.Join(configDir, "config.yaml")

	// Default config
	cfg := &AppConfig{
		Paths: PathConfig{
			VMs:   filepath.Join(GetDefaultDataDir(), "machines"),
			Cache: GetDefaultCacheDir(),
		},
		QEMU: QEMUConfig{
			Binary: "auto",
		},
		Server: ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
		},
	}

	// Try to load from file
	if data, err := os.ReadFile(cfgPath); err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	// ENV overrides
	if val := os.Getenv("VMTOOL_HOST"); val != "" {
		cfg.Server.Host = val
	}
	if val := os.Getenv("VMTOOL_PORT"); val != "" {
		var port int
		if _, err := fmt.Sscanf(val, "%d", &port); err != nil {
			return nil, fmt.Errorf("invalid VMTOOL_PORT: %v", err)
		}
		cfg.Server.Port = port
	}

	return cfg, nil
}
