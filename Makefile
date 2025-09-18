.PHONY: all clean server wasm web test fmt vet lint deps help

# Default target
all: server wasm

# Build server binary
server:
	@echo "Building server..."
	go build -o bin/server ./cmd/server

# Build WASM binary
wasm:
	@echo "Building WASM..."
	GOOS=js GOARCH=wasm go build -o web/main.wasm ./cmd/wasm

# Copy wasm_exec.js from Go installation
web: wasm
	@echo "Setting up web directory..."
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" web/

# Run server (builds if necessary)
run-server: server
	@echo "Starting server..."
	./bin/server

# Run development server with file watching (requires fswatch)
dev: web
	@echo "Starting development server with auto-rebuild..."
	@if command -v fswatch >/dev/null 2>&1; then \
		fswatch -o . -e ".*" -i "\\.go$$" | xargs -n1 -I{} make web & \
	fi
	make run-server

# Test all packages
test:
	@echo "Running tests..."
	go test -v ./...

# Test with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -cover ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	go vet ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Create necessary directories
dirs:
	@echo "Creating directories..."
	mkdir -p bin web

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -f web/main.wasm
	rm -f web/wasm_exec.js
	go clean

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@if command -v brew >/dev/null 2>&1; then \
		brew install fswatch; \
	elif command -v apt-get >/dev/null 2>&1; then \
		sudo apt-get install inotify-tools; \
	else \
		echo "Please install fswatch manually for file watching support"; \
	fi

# Help target
help:
	@echo "Available targets:"
	@echo "  all          - Build both server and WASM"
	@echo "  server       - Build server binary"
	@echo "  wasm         - Build WASM binary"
	@echo "  web          - Build WASM and copy wasm_exec.js"
	@echo "  run-server   - Build and run server"
	@echo "  dev          - Start development server with auto-rebuild"
	@echo "  test         - Run all tests"
	@echo "  test-coverage- Run tests with coverage"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  lint         - Run linter"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  dirs         - Create necessary directories"
	@echo "  clean        - Clean build artifacts"
	@echo "  install-tools- Install development tools"
	@echo "  help         - Show this help"

# Create directories before building
server wasm: dirs