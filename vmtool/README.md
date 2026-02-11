# VMTool - Terminal VM Streaming

VMTool is a terminal-based alternative to UTM that provisions virtual machines and streams them to a browser. It is built in Go and uses QEMU as the virtualization engine, providing the same powerful virtualization capabilities as the UTM desktop app in a CLI-first interface.

## Features

- **CLI-First**: Manage VMs directly from your terminal.
- **Web Display**: Stream VM output to any browser via noVNC (embedded).
- **Universal Architecture Support**: Create and run VMs for x86_64, ARM64 (aarch64), and other architectures on any host platform.
- **Hardware Acceleration**: Automatic detection and use of platform-specific accelerators:
  - macOS: Hypervisor.framework (HVF) for both Intel and Apple Silicon
  - Linux: KVM for Intel/AMD processors
  - Windows: Windows Hypervisor Platform (WHPX)
  - Fallback: TCG software emulation for any architecture on any platform
- **Cross-Platform Host**: Runs on macOS, Linux, and Windows.
- **UTM Import**: Easily import existing `.utm` bundles from UTM desktop app.
- **REST API**: Control VMs programmatically.
- **Secure**: Token-based authentication and localhost binding by default.

## Installation

### Prerequisites

- QEMU 7.0 or newer must be installed on your system.

### Quick Install (Recommended)

Use the installation script to automatically download and install vmtool:

```bash
curl -fsSL https://raw.githubusercontent.com/badbsallyy/UTM/main/vmtool/install.sh | sudo bash
```

Or if you have the repository cloned:

```bash
cd vmtool
sudo ./install.sh
```

### Manual Installation

#### Option 1: Install from Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/badbsallyy/UTM/releases):

```bash
# For Linux (x86_64)
curl -L -o vmtool https://github.com/badbsallyy/UTM/releases/latest/download/vmtool-linux-amd64
chmod +x vmtool
sudo mv vmtool /usr/local/bin/

# For macOS (Apple Silicon)
curl -L -o vmtool https://github.com/badbsallyy/UTM/releases/latest/download/vmtool-darwin-arm64
chmod +x vmtool
sudo mv vmtool /usr/local/bin/

# For macOS (Intel)
curl -L -o vmtool https://github.com/badbsallyy/UTM/releases/latest/download/vmtool-darwin-amd64
chmod +x vmtool
sudo mv vmtool /usr/local/bin/
```

#### Option 2: Build from Source

Clone the repository and build:

```bash
git clone https://github.com/badbsallyy/UTM.git
cd UTM/vmtool
make build
```

The binary will be located in `build/vmtool`.

To install system-wide:

```bash
# Install to /usr/local/bin (requires sudo)
sudo make install-system

# Or install to your GOPATH/bin (no sudo required)
make install
```

### Verify Installation

After installation, verify that vmtool is accessible from your terminal:

```bash
vmtool --help
```

## Usage

### Initialize

```bash
vmtool init
```

### Create a VM

```bash
vmtool create my-ubuntu
```

### Import a UTM bundle

```bash
vmtool import ~/Downloads/Ubuntu.utm
```

### Start the server

```bash
vmtool serve
```

Access the dashboard at `http://localhost:8080`.

### Control from CLI

```bash
vmtool start my-ubuntu
vmtool stop my-ubuntu
vmtool pause my-ubuntu
vmtool resume my-ubuntu
vmtool list
vmtool info my-ubuntu
```

### Snapshots

```bash
vmtool snapshot create my-ubuntu snap1
vmtool snapshot restore my-ubuntu snap1
vmtool snapshot delete my-ubuntu snap1
```

## Architecture

VMTool leverages QEMU to provide universal VM creation capabilities across different host platforms:

- **Backend**: Go 1.21+
- **Virtualization**: QEMU + QMP (QEMU Machine Protocol)
- **Streaming**: VNC + WebSocket Proxy + noVNC (HTML5 Canvas)
- **Configuration**: YAML-based VM configurations

### Supported Guest Architectures

VMTool can create and run VMs for any architecture supported by QEMU, including:
- x86_64 (Intel/AMD 64-bit)
- i386 (Intel/AMD 32-bit)
- aarch64 (ARM 64-bit / Apple Silicon)
- armv7 (ARM 32-bit)
- riscv64, ppc64, sparc64, and more

### Platform-Specific Acceleration

| Host Platform        | Accelerator | Use Case                                  |
|----------------------|-------------|-------------------------------------------|
| macOS (Intel)        | HVF         | Fast x86_64 guest VMs                     |
| macOS (Apple Silicon)| HVF         | Fast aarch64 guest VMs                    |
| Linux (x86_64)       | KVM         | Fast x86_64 guest VMs                     |
| Linux (ARM64)        | KVM         | Fast aarch64 guest VMs                    |
| Windows              | WHPX        | Fast x86_64 guest VMs                     |
| Any Platform         | TCG         | Software emulation for any guest architecture |

This enables true universal VM creation: run x86_64 VMs on Apple Silicon, ARM VMs on Intel, or any architecture on any platform using the same QEMU backend as the UTM desktop application.
