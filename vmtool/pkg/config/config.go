package config

import (
	"os"
	"path/filepath"
	"runtime"
)

type VMConfig struct {
	Name           string         `yaml:"name"`
	UUID           string         `yaml:"uuid"`
	System         SystemConfig   `yaml:"system"`
	Drives         []DriveConfig  `yaml:"drives"`
	Network        NetworkConfig  `yaml:"network"`
	Display        DisplayConfig  `yaml:"display"`
	Sharing        SharingConfig  `yaml:"sharing"`
	Boot           BootConfig     `yaml:"boot"`
	AdditionalArgs []string       `yaml:"additional_args,omitempty"`
}

type SystemConfig struct {
	Architecture string `yaml:"architecture"` // e.g., x86_64, aarch64
	Target       string `yaml:"target"`       // e.g., pc-q35-7.0, virt
	CPU          string `yaml:"cpu"`          // e.g., host, Skylake-Client
	Memory       int    `yaml:"memory"`       // MB
	CPUs         int    `yaml:"cpus"`         // Count
	Accelerator  string `yaml:"accelerator"`  // hvf, kvm, whpx, tcg
}

type BootConfig struct {
	Order []string `yaml:"order"` // e.g., ["disk", "cdrom", "network"]
}

type DriveConfig struct {
	ID        int    `yaml:"id"`
	Interface string `yaml:"interface"` // ide, scsi, virtio, nvme, usb
	ImagePath string `yaml:"image_path"`
	ImageType string `yaml:"image_type"` // disk, cdrom, bios, kernel, initrd
	ReadOnly  bool   `yaml:"read_only"`
}

type NetworkConfig struct {
	Mode           string           `yaml:"mode"` // user (slirp), bridged
	Hardware       string           `yaml:"hardware"`
	PortForwards   []PortForward    `yaml:"port_forwards,omitempty"`
}

type PortForward struct {
	Protocol string `yaml:"protocol"` // tcp, udp
	HostPort int    `yaml:"host_port"`
	GuestPort int    `yaml:"guest_port"`
	HostIP   string `yaml:"host_ip,omitempty"`
	GuestIP  string `yaml:"guest_ip,omitempty"`
}

type DisplayConfig struct {
	Enabled bool   `yaml:"enabled"`
	VNCAddr string `yaml:"vnc_addr"` // e.g., :0 or localhost:5901
	VNCPort int    `yaml:"vnc_port,omitempty"` // e.g., 5900
	Width   int    `yaml:"width,omitempty"`
	Height  int    `yaml:"height,omitempty"`
}

type AppConfig struct {
	Paths    PathConfig    `yaml:"paths"`
	QEMU     QEMUConfig     `yaml:"qemu"`
	Server   ServerConfig   `yaml:"server"`
	Security SecurityConfig `yaml:"security"`
}

type PathConfig struct {
	VMs   string `yaml:"vms"`
	Cache string `yaml:"cache"`
}

type QEMUConfig struct {
	Binary string `yaml:"binary"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type SecurityConfig struct {
	APIToken string `yaml:"api_token"`
}

type SharingConfig struct {
	DirectoryShare string `yaml:"directory_share,omitempty"`
	ReadOnly       bool   `yaml:"read_only,omitempty"`
}

func GetDefaultDataDir() string {
	if val := os.Getenv("VMTOOL_HOME"); val != "" {
		return val
	}
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "vmtool")
	case "windows":
		return filepath.Join(os.Getenv("LOCALAPPDATA"), "vmtool")
	default: // linux
		return filepath.Join(home, ".local", "share", "vmtool")
	}
}

func GetDefaultConfigDir() string {
	if val := os.Getenv("VMTOOL_HOME"); val != "" {
		return val
	}
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "vmtool")
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "vmtool")
	default: // linux
		return filepath.Join(home, ".config", "vmtool")
	}
}

func GetDefaultCacheDir() string {
	if val := os.Getenv("VMTOOL_HOME"); val != "" {
		return filepath.Join(val, "cache")
	}
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Caches", "vmtool")
	case "windows":
		return filepath.Join(os.Getenv("LOCALAPPDATA"), "vmtool", "cache")
	default: // linux
		return filepath.Join(home, ".cache", "vmtool")
	}
}
