.PHONY: build clean test install install-system uninstall help cross-build

# Build variables
BINARY_NAME=vmtool
BUILD_DIR=build
GO=go
GOFLAGS=-v
PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -s -w -X main.version=$(VERSION)

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	$(GO) test ./... -v

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME) $(BINARY_NAME).exe vmtool-test
	@echo "Clean complete"

# Install to GOPATH/bin (user installation)
install: build
	@echo "Installing $(BINARY_NAME) to GOPATH/bin..."
	$(GO) install .
	@echo "Install complete"

# Install to system location (requires sudo)
install-system: build
	@echo "Installing $(BINARY_NAME) to $(BINDIR)..."
	@install -d $(BINDIR)
	@install -m 755 $(BUILD_DIR)/$(BINARY_NAME) $(BINDIR)/$(BINARY_NAME)
	@echo "Install complete: $(BINDIR)/$(BINARY_NAME)"
	@echo "You can now use '$(BINARY_NAME)' from anywhere in your terminal"

# Uninstall from system
uninstall:
	@echo "Uninstalling $(BINARY_NAME) from $(BINDIR)..."
	@rm -f $(BINDIR)/$(BINARY_NAME)
	@echo "Uninstall complete"

# Cross-compile for all platforms
cross-build:
	@echo "Cross-compiling for all platforms..."
	@mkdir -p $(BUILD_DIR)
	# Linux
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	# macOS
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	# Windows
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "Cross-compilation complete"

# Display help
help:
	@echo "Available targets:"
	@echo "  build           - Build the vmtool binary (default)"
	@echo "  test            - Run all tests"
	@echo "  clean           - Remove build artifacts"
	@echo "  install         - Build and install to GOPATH/bin"
	@echo "  install-system  - Build and install to $(BINDIR) (requires sudo)"
	@echo "  uninstall       - Remove vmtool from $(BINDIR)"
	@echo "  cross-build     - Build for all supported platforms"
	@echo "  help            - Display this help message"
	@echo ""
	@echo "Installation options:"
	@echo "  make install         - Install to GOPATH/bin (no sudo required)"
	@echo "  make install-system  - Install to /usr/local/bin (requires sudo)"
	@echo "  ./install.sh         - Run interactive installation script"
