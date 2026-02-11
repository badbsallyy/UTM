# VMTool - Terminal VM Streaming

VMTool is a terminal-based alternative to UTM that provisions virtual machines and streams them to a browser. It is built in Go and uses QEMU as the virtualization engine.

## Features

- **CLI-First**: Manage VMs directly from your terminal.
- **Web Display**: Stream VM output to any browser via noVNC (embedded).
- **Cross-Platform**: Supports macOS (HVF), Linux (KVM), and Windows (WHPX).
- **UTM Import**: Easily import existing `.utm` bundles.
- **REST API**: Control VMs programmatically.
- **Secure**: Token-based authentication and localhost binding by default.

## Installation

### Prerequisites

- QEMU 7.0 or newer must be installed on your system.

### Build from source

```bash
make build
```

The binary will be located in `build/vmtool`.

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

- **Backend**: Go 1.21+
- **Virtualization**: QEMU + QMP
- **Streaming**: VNC + WebSocket Proxy + noVNC (HTML5 Canvas)
- **Configuration**: YAML
