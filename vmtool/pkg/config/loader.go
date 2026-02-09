package config

import (
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
		yaml.Unmarshal(data, cfg)
	}

	// ENV overrides
	if val := os.Getenv("VMTOOL_HOST"); val != "" {
		cfg.Server.Host = val
	}
	if val := os.Getenv("VMTOOL_PORT"); val != "" {
		// handle port conversion if needed
	}

	return cfg, nil
}
