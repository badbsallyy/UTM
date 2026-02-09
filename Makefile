BINARY_NAME=vmtool
BUILD_DIR=build

all: build

build:
	mkdir -p $(BUILD_DIR)
	cd vmtool && go build -o ../$(BUILD_DIR)/$(BINARY_NAME) .

clean:
	rm -rf $(BUILD_DIR)

test:
	cd vmtool && go test ./...

cross-build:
	mkdir -p $(BUILD_DIR)
	# Linux
	GOOS=linux GOARCH=amd64 cd vmtool && go build -o ../$(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 cd vmtool && go build -o ../$(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	# macOS
	GOOS=darwin GOARCH=amd64 cd vmtool && go build -o ../$(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 cd vmtool && go build -o ../$(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	# Windows
	GOOS=windows GOARCH=amd64 cd vmtool && go build -o ../$(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

.PHONY: all build clean test cross-build
