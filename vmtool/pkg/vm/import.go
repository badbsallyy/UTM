package vm

import (
	"os"
	"path/filepath"

	"howett.net/plist"
	"github.com/utmapp/vmtool/pkg/config"
)

type UTMConfig struct {
	Information struct {
		Name string `plist:"Name"`
		UUID string `plist:"UUID"`
	} `plist:"Information"`
	System struct {
		Architecture string   `plist:"Architecture"`
		Memory       int      `plist:"Memory"`
		CPUCount     int      `plist:"CPUCount"`
		Target       string   `plist:"Target"`
		BootOrder    []string `plist:"BootOrder"`
	} `plist:"System"`
	Drives []struct {
		ImageName string `plist:"ImageName"`
		Interface string `plist:"Interface"`
		ImageType string `plist:"ImageType"`
	} `plist:"Drive"`
	Networks []struct {
		Hardware   string `plist:"Hardware"`
		NetworkMode string `plist:"NetworkMode"`
		MACAddress string `plist:"MACAddress"`
	} `plist:"Network"`
	Displays []struct {
		Upscaling   string `plist:"Upscaling"`
		Downscaling string `plist:"Downscaling"`
	} `plist:"Display"`
}

func ImportUTM(bundlePath string) (*config.VMConfig, []string, error) {
	var warnings []string
	configPath := filepath.Join(bundlePath, "config.plist")
	f, err := os.Open(configPath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	var utmCfg UTMConfig
	decoder := plist.NewDecoder(f)
	if err := decoder.Decode(&utmCfg); err != nil {
		return nil, nil, err
	}

	vmCfg := &config.VMConfig{
		Name: utmCfg.Information.Name,
		UUID: utmCfg.Information.UUID,
		System: config.SystemConfig{
			Architecture: utmCfg.System.Architecture,
			Memory:       utmCfg.System.Memory,
			CPUs:         utmCfg.System.CPUCount,
			Target:       utmCfg.System.Target,
		},
		Display: config.DisplayConfig{
			Enabled: true,
			VNCAddr: ":0",
		},
		Boot: config.BootConfig{
			Order: utmCfg.System.BootOrder,
		},
	}

	for i, d := range utmCfg.Drives {
		imagePath := filepath.Join(bundlePath, "Data", d.ImageName)
		vmCfg.Drives = append(vmCfg.Drives, config.DriveConfig{
			ID:        i,
			Interface: d.Interface,
			ImagePath: imagePath,
			ImageType: d.ImageType,
		})
	}

	if len(utmCfg.Networks) > 0 {
		n := utmCfg.Networks[0]
		vmCfg.Network = config.NetworkConfig{
			Mode:     n.NetworkMode,
			Hardware: n.Hardware,
		}
		if n.NetworkMode == "shared" {
			vmCfg.Network.Mode = "user"
		}
		// MAC address could be added to NetworkConfig if needed
	}

	warnings = append(warnings, "Converted SPICE to VNC (SPICE not supported)")

	return vmCfg, warnings, nil
}
