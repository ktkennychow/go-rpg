.PHONY: dev build run clean install-deps

# Development with hot reload
# Check dependencies, build WASM binary, start live-server in background, watch for Go file changes
dev:
	@echo "Starting development server..."
	@if ! command -v watchexec >/dev/null 2>&1; then \
		echo "Error: watchexec not found. Install with: brew install watchexec"; \
		exit 1; \
	fi
	@if ! command -v live-server >/dev/null 2>&1; then \
		echo "Error: live-server not found. Install with: npm install -g live-server"; \
		exit 1; \
	fi
	@echo "Building initial WASM binary..."
	@env GOOS=js GOARCH=wasm go build -o main.wasm github.com/ktkennychow/go-rpg
	@echo "Starting live-server in background..."
	@live-server . & echo $$! > /tmp/go-rpg-live-server.pid
	@echo "Watching for Go file changes..."
	@watchexec -e go -r "env GOOS=js GOARCH=wasm go build -o main.wasm github.com/ktkennychow/go-rpg" || (kill `cat /tmp/go-rpg-live-server.pid` 2>/dev/null; rm -f /tmp/go-rpg-live-server.pid)

# Build for local development
build:
	@go build -o go-rpg-native .

# Build WASM
build-wasm:
	@env GOOS=js GOARCH=wasm go build -o main.wasm github.com/ktkennychow/go-rpg

# Run locally (native)
run:
	@go run .

# Clean build artifacts
clean:
	@rm -f go-rpg-native main.wasm /tmp/go-rpg-live-server.pid
	@pkill -f live-server || true

# Install development dependencies
install-deps:
	@echo "Installing development dependencies..."
	@if command -v brew >/dev/null 2>&1; then \
		echo "Installing watchexec via brew..."; \
		brew install watchexec; \
	else \
		echo "Brew not found. Please install watchexec manually."; \
	fi
	@if command -v npm >/dev/null 2>&1; then \
		echo "Installing live-server via npm..."; \
		npm install -g live-server; \
	else \
		echo "npm not found. Please install live-server manually."; \
	fi

# Help
help:
	@echo "Available commands:"
	@echo "  make dev          - Start development server with hot reload"
	@echo "  make build        - Build native binary"
	@echo "  make build-wasm   - Build WASM binary"
	@echo "  make run          - Run locally (native)"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make install-deps - Install development dependencies"
	@echo "  make help         - Show this help"
