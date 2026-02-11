package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/cast"
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
		cfg.Server.Port = cast.ToInt(val)
	}

	return cfg, nil
}
