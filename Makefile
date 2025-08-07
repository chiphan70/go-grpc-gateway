# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=go-grpcgateway
BINARY_UNIX=$(BINARY_NAME)_unix

# Proto parameters
PROTO_DIR=api/proto
PB_DIR=pkg/pb

# Docker parameters
DOCKER_IMAGE=go-grpcgateway
DOCKER_TAG=latest

.PHONY: all build clean test deps proto help

# Default target
all: deps buf-generate build

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Install required tools
tools:
	@echo "Installing required tools..."
	$(GOGET) google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GOGET) google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	$(GOGET) github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	$(GOGET) github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# Install Buf CLI
install-buf:
	@echo "Installing Buf CLI..."
	@if ! command -v buf &> /dev/null; then \
		echo "Installing Buf CLI..."; \
		if [[ "$$OSTYPE" == "darwin"* ]]; then \
			brew install bufbuild/buf/buf; \
		else \
			curl -sSL "https://github.com/bufbuild/buf/releases/latest/download/buf-$$(uname -s)-$$(uname -m)" -o "/tmp/buf" && \
			chmod +x "/tmp/buf" && \
			sudo mv "/tmp/buf" "/usr/local/bin/buf"; \
		fi; \
	else \
		echo "Buf CLI is already installed: $$(buf --version)"; \
	fi

# Generate protobuf files using protoc (legacy)
proto:
	@echo "Generating protobuf files with protoc..."
	@mkdir -p $(PB_DIR)
	protoc \
		-I $(PROTO_DIR) \
		-I $(shell go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway/v2)/third_party/googleapis \
		--go_out=paths=source_relative:$(PB_DIR) \
		--go-grpc_out=paths=source_relative:$(PB_DIR) \
		--grpc-gateway_out=paths=source_relative:$(PB_DIR) \
		--openapiv2_out=. \
		$(PROTO_DIR)/*.proto

# Generate protobuf files using Buf CLI (recommended)
buf-generate:
	@echo "Generating protobuf files with Buf..."
	@mkdir -p $(PB_DIR) docs
	buf generate

# Lint protobuf files
buf-lint:
	@echo "Linting protobuf files..."
	buf lint

# Format protobuf files
buf-format:
	@echo "Formatting protobuf files..."
	buf format -w

# Check for breaking changes
buf-breaking:
	@echo "Checking for breaking changes..."
	buf breaking --against '.git#branch=main'

# Update Buf dependencies
buf-deps:
	@echo "Updating Buf dependencies..."
	buf mod update

# Build the binary
build:
	@echo "Building..."
	$(GOBUILD) -o bin/$(BINARY_NAME) -v ./cmd

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_UNIX) -v ./cmd

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf bin/
	rm -rf $(PB_DIR)/
	rm -rf docs/

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run the application
run:
	@echo "Running application..."
	$(GOBUILD) -o bin/$(BINARY_NAME) -v ./cmd && ./bin/$(BINARY_NAME)

# Run with live reload (requires air)
dev:
	@echo "Running with live reload..."
	air

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	golangci-lint run

# Check for security issues (requires gosec)
security:
	@echo "Checking for security issues..."
	gosec ./...

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 -p 9090:9090 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

# Docker compose up
docker-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

# Docker compose down
docker-down:
	@echo "Stopping services with Docker Compose..."
	docker-compose down

# Show help
help:
	@echo "Available targets:"
	@echo "  all           - Download deps, generate proto with Buf, and build"
	@echo "  deps          - Download dependencies"
	@echo "  tools         - Install required protoc tools"
	@echo "  install-buf   - Install Buf CLI"
	@echo "  proto         - Generate protobuf files (legacy protoc)"
	@echo "  buf-generate  - Generate protobuf files with Buf (recommended)"
	@echo "  buf-lint      - Lint protobuf files"
	@echo "  buf-format    - Format protobuf files"
	@echo "  buf-breaking  - Check for breaking changes"
	@echo "  buf-deps      - Update Buf dependencies"
	@echo "  build         - Build the binary"
	@echo "  build-linux   - Build for Linux"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  run           - Build and run the application"
	@echo "  dev           - Run with live reload"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  security      - Check for security issues"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  docker-up     - Start services with Docker Compose"
	@echo "  docker-down   - Stop services with Docker Compose"
	@echo "  help          - Show this help"
