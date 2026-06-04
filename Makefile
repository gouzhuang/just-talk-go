.PHONY: build build-gnome install install-gnome test clean run run-gnome deps

APP_NAME := just-talk
CMD_DIR := ./cmd/just-talk
BUILD_DIR := ./build

# Build for current platform (default frontend)
build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

# Build with GNOME system tray support
build-gnome:
	go build -tags gnome -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

# Install to ~/.local/bin
install: build
	$(BUILD_DIR)/$(APP_NAME) --install

# Install GNOME build
install-gnome: build-gnome
	$(BUILD_DIR)/$(APP_NAME) --install

# Run (default frontend)
run:
	go run $(CMD_DIR)

# Run with GNOME frontend
run-gnome:
	go run -tags gnome $(CMD_DIR) --frontend gnome

# Test
test:
	go test ./... -v

# Clean
clean:
	rm -rf $(BUILD_DIR)

# Install dependencies
deps:
	go mod tidy
	go mod download
