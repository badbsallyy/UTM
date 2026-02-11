#!/bin/bash
# VMTool Installation Script
# This script downloads and installs vmtool to your system

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Detect OS and architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case "$os" in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        mingw*|msys*|cygwin*)
            OS="windows"
            ;;
        *)
            echo -e "${RED}Unsupported operating system: $os${NC}"
            exit 1
            ;;
    esac
    
    case "$arch" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            echo -e "${RED}Unsupported architecture: $arch${NC}"
            exit 1
            ;;
    esac
    
    BINARY_NAME="vmtool-${OS}-${ARCH}"
    if [ "$OS" = "windows" ]; then
        BINARY_NAME="${BINARY_NAME}.exe"
    fi
}

# Check if running with sudo (for system-wide install)
check_sudo() {
    if [ "$EUID" -ne 0 ] && [ "$INSTALL_DIR" = "/usr/local/bin" ]; then
        echo -e "${YELLOW}Installing to /usr/local/bin requires sudo privileges.${NC}"
        echo "Please run with sudo: sudo $0"
        exit 1
    fi
}

# Install from local build
install_local() {
    local script_dir="$(cd "$(dirname "$0")" && pwd)"
    local build_dir="${script_dir}/build"
    local binary_path="${build_dir}/vmtool"
    
    if [ ! -f "$binary_path" ]; then
        echo -e "${YELLOW}Binary not found at $binary_path${NC}"
        echo "Building vmtool..."
        cd "$script_dir"
        make build
    fi
    
    echo -e "${GREEN}Installing vmtool from local build...${NC}"
    install -m 755 "$binary_path" "${INSTALL_DIR}/vmtool"
    echo -e "${GREEN}✓ vmtool installed to ${INSTALL_DIR}/vmtool${NC}"
}

# Install from GitHub releases
install_from_github() {
    local latest_release=$(curl -s https://api.github.com/repos/badbsallyy/UTM/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$latest_release" ]; then
        echo -e "${RED}Could not fetch latest release information${NC}"
        echo "Falling back to local build..."
        install_local
        return
    fi
    
    local download_url="https://github.com/badbsallyy/UTM/releases/download/${latest_release}/${BINARY_NAME}"
    
    echo -e "${GREEN}Downloading vmtool ${latest_release} for ${OS}/${ARCH}...${NC}"
    
    if command -v curl &> /dev/null; then
        curl -L -o "/tmp/vmtool" "$download_url"
    elif command -v wget &> /dev/null; then
        wget -O "/tmp/vmtool" "$download_url"
    else
        echo -e "${RED}Neither curl nor wget found. Please install one of them.${NC}"
        exit 1
    fi
    
    if [ ! -f "/tmp/vmtool" ]; then
        echo -e "${RED}Download failed${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}Installing vmtool...${NC}"
    install -m 755 "/tmp/vmtool" "${INSTALL_DIR}/vmtool"
    rm -f "/tmp/vmtool"
    echo -e "${GREEN}✓ vmtool installed to ${INSTALL_DIR}/vmtool${NC}"
}

# Main installation function
main() {
    echo "VMTool Installation Script"
    echo "=========================="
    echo ""
    
    detect_platform
    
    # Set installation directory
    INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
    
    # Check if we need sudo
    check_sudo
    
    # Check if we're in the vmtool directory
    if [ -f "Makefile" ] && [ -f "main.go" ]; then
        # We're in the source directory, install from local build
        install_local
    else
        # Not in source directory, try to download from GitHub
        echo -e "${YELLOW}Not in vmtool source directory, attempting to download from GitHub...${NC}"
        install_from_github
    fi
    
    # Verify installation
    if command -v vmtool &> /dev/null; then
        echo ""
        echo -e "${GREEN}Installation successful!${NC}"
        echo ""
        vmtool --version || echo "vmtool is installed at: $(which vmtool)"
        echo ""
        echo "You can now use vmtool from anywhere in your terminal."
        echo "Try: vmtool --help"
    else
        echo ""
        echo -e "${YELLOW}Installation completed, but vmtool is not in PATH.${NC}"
        echo "Please add ${INSTALL_DIR} to your PATH, or use the full path: ${INSTALL_DIR}/vmtool"
    fi
}

# Run main function
main "$@"
