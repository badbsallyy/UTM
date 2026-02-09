package qemu

import (
	"fmt"
	"runtime"

	"github.com/utmapp/vmtool/pkg/config"
)

type Builder struct {
	config *config.VMConfig
}

func NewBuilder(cfg *config.VMConfig) *Builder {
	return &Builder{config: cfg}
}

func (b *Builder) BuildArgs() []string {
	var args []string

	// Basic identity
	args = append(args, "-name", b.config.Name)
	args = append(args, "-uuid", b.config.UUID)

	// System settings
	args = append(args, "-m", fmt.Sprintf("%d", b.config.System.Memory))
	args = append(args, "-smp", fmt.Sprintf("cpus=%d", b.config.System.CPUs))

	// Accelerator
	accel := b.config.System.Accelerator
	if accel == "" {
		accel = b.DetectAccelerator()
	}
	args = append(args, "-accel", accel)

	// Machine and CPU
	if b.config.System.Target != "" {
		args = append(args, "-machine", b.config.System.Target)
	}
	if b.config.System.CPU != "" {
		args = append(args, "-cpu", b.config.System.CPU)
	}

	// Drives
	for _, drive := range b.config.Drives {
		args = append(args, b.buildDriveArgs(drive)...)
	}

	// Network
	args = append(args, b.buildNetworkArgs()...)

	// Display (VNC)
	if b.config.Display.Enabled {
		vncAddr := b.config.Display.VNCAddr
		if b.config.Display.VNCPort != 0 {
			vncAddr = fmt.Sprintf("127.0.0.1:%d", b.config.Display.VNCPort-5900)
		}
		args = append(args, "-vnc", vncAddr)
	} else {
		args = append(args, "-nographic")
	}

	// Additional arguments
	args = append(args, b.config.AdditionalArgs...)

	return args
}

func (b *Builder) DetectAccelerator() string {
	switch runtime.GOOS {
	case "darwin":
		return "hvf"
	case "linux":
		return "kvm"
	case "windows":
		return "whpx"
	default:
		return "tcg"
	}
}

func (b *Builder) buildDriveArgs(drive config.DriveConfig) []string {
	var args []string
	driveID := fmt.Sprintf("drive%d", drive.ID)

	// -drive if=none,id=drive0,file=...,format=qcow2
	fileArg := fmt.Sprintf("if=none,id=%s,file=%s", driveID, drive.ImagePath)
	if drive.ReadOnly {
		fileArg += ",readonly=on"
	}
	args = append(args, "-drive", fileArg)

	// -device virtio-blk-pci,drive=drive0
	deviceType := b.getDeviceType(drive.Interface, drive.ImageType)
	args = append(args, "-device", fmt.Sprintf("%s,drive=%s", deviceType, driveID))

	return args
}

func (b *Builder) getDeviceType(iface, imageType string) string {
	if imageType == "cdrom" {
		return "ide-cd" // Simplified
	}
	switch iface {
	case "virtio":
		return "virtio-blk-pci"
	case "nvme":
		return "nvme"
	case "usb":
		return "usb-storage"
	default:
		return "ide-hd"
	}
}

func (b *Builder) buildNetworkArgs() []string {
	var args []string
	// Default to user mode slirp
	netdev := "user,id=net0"
	for _, fw := range b.config.Network.PortForwards {
		netdev += fmt.Sprintf(",hostfwd=%s::%d-:%d", fw.Protocol, fw.HostPort, fw.GuestPort)
	}
	args = append(args, "-netdev", netdev)

	hardware := b.config.Network.Hardware
	if hardware == "" {
		hardware = "virtio-net-pci"
	}
	args = append(args, "-device", fmt.Sprintf("%s,netdev=net0", hardware))

	return args
}
