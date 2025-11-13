# RESTerX CLI Build Makefile

# Binary name
BINARY_NAME=resterx-cli

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.CommitHash=$(COMMIT_HASH) -X main.BuildDate=$(BUILD_DATE)"

# Build directory
DIST_DIR=dist
BUILD_DIR=build

# Platforms
PLATFORMS=windows/amd64 darwin/amd64 darwin/arm64 linux/amd64

.PHONY: all build clean test help install deps build-all

# Default target
all: clean build-all checksums

# Help target
help:
	@echo "RESTerX CLI Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  make build         - Build for current platform"
	@echo "  make build-all     - Build for all platforms"
	@echo "  make checksums     - Generate SHA256 checksums"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make deps          - Download dependencies"
	@echo "  make install       - Install binary locally"
	@echo ""
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT_HASH)"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Build for current platform
build: deps
	@echo "Building $(BINARY_NAME) for current platform..."
	@mkdir -p $(DIST_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME) ./cmd
	@echo "Build complete: $(DIST_DIR)/$(BINARY_NAME)"

# Build for all platforms
build-all: deps
	@echo "Building $(BINARY_NAME) for all platforms..."
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		output_name=$(DIST_DIR)/$(BINARY_NAME)-$${GOOS}-$${GOARCH}; \
		if [ "$$GOOS" = "windows" ]; then output_name=$${output_name}.exe; fi; \
		echo "Building for $$GOOS/$$GOARCH..."; \
		GOOS=$$GOOS GOARCH=$$GOARCH CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $$output_name ./cmd; \
		if [ $$? -ne 0 ]; then \
			echo "Error building for $$GOOS/$$GOARCH"; \
			exit 1; \
		fi; \
	done
	@echo "All builds complete!"
	@ls -lh $(DIST_DIR)

# Generate SHA256 checksums
checksums:
	@echo "Generating SHA256 checksums..."
	@cd $(DIST_DIR) && \
		for file in $(BINARY_NAME)-*; do \
			if [ -f "$$file" ]; then \
				sha256sum "$$file" >> checksums.txt; \
			fi; \
		done
	@echo "Checksums saved to $(DIST_DIR)/checksums.txt"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(DIST_DIR)
	@rm -rf $(BUILD_DIR)
	$(GOCLEAN)
	@echo "Clean complete!"

# Install binary locally
install: build
	@echo "Installing $(BINARY_NAME) to GOPATH/bin..."
	@cp $(DIST_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME) 2>/dev/null || cp $(DIST_DIR)/$(BINARY_NAME) ~/go/bin/$(BINARY_NAME)
	@echo "Installation complete!"
