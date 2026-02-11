.PHONY: build clean test install help

# Build variables
BINARY_NAME=vmtool
BUILD_DIR=build
GO=go
GOFLAGS=-v

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
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

# Install to system
install: build
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install .
	@echo "Install complete"

# Display help
help:
	@echo "Available targets:"
	@echo "  build    - Build the vmtool binary (default)"
	@echo "  test     - Run all tests"
	@echo "  clean    - Remove build artifacts"
	@echo "  install  - Build and install to GOPATH/bin"
	@echo "  help     - Display this help message"
