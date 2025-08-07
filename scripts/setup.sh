#!/bin/bash

# Setup script for go-grpcgateway project

set -e

echo "🚀 Setting up go-grpcgateway project..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

echo "✅ Go is installed: $(go version)"

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "❌ protoc is not installed. Please install Protocol Buffers compiler."
    echo "For macOS: brew install protobuf"
    echo "For Ubuntu: sudo apt-get install protobuf-compiler"
    exit 1
fi

echo "✅ protoc is installed: $(protoc --version)"

# Install Buf CLI
echo "📦 Installing Buf CLI..."
if ! command -v buf &> /dev/null; then
    echo "Installing Buf CLI..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        if command -v brew &> /dev/null; then
            brew install bufbuild/buf/buf
        else
            echo "❌ Homebrew not found. Installing Buf manually..."
            curl -sSL "https://github.com/bufbuild/buf/releases/latest/download/buf-$(uname -s)-$(uname -m)" -o "/tmp/buf"
            chmod +x "/tmp/buf"
            sudo mv "/tmp/buf" "/usr/local/bin/buf"
        fi
    else
        curl -sSL "https://github.com/bufbuild/buf/releases/latest/download/buf-$(uname -s)-$(uname -m)" -o "/tmp/buf"
        chmod +x "/tmp/buf"
        sudo mv "/tmp/buf" "/usr/local/bin/buf"
    fi
else
    echo "✅ Buf CLI is already installed: $(buf --version)"
fi

# Install Go tools (optional, only needed if using protoc directly)
echo "📦 Installing Go protoc tools (optional)..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# Download dependencies
echo "📦 Downloading dependencies..."
go mod download
go mod tidy

# Create necessary directories
echo "📁 Creating directories..."
mkdir -p bin
mkdir -p pkg/pb

# Generate protobuf files with Buf
echo "🔧 Generating protobuf files with Buf..."
make buf-generate

# Build the project
echo "🔨 Building the project..."
make build

# Check if MongoDB is running (optional)
if command -v mongosh &> /dev/null; then
    if mongosh --eval "db.runCommand('ping')" &> /dev/null; then
        echo "✅ MongoDB is running"
    else
        echo "⚠️  MongoDB is not running. You can start it with: brew services start mongodb/brew/mongodb-community"
    fi
elif command -v mongo &> /dev/null; then
    if mongo --eval "db.runCommand('ping')" &> /dev/null; then
        echo "✅ MongoDB is running"
    else
        echo "⚠️  MongoDB is not running. You can start it with: brew services start mongodb/brew/mongodb-community"
    fi
else
    echo "⚠️  MongoDB client not found. Please install MongoDB."
fi

echo ""
echo "🎉 Setup completed successfully!"
echo ""
echo "To run the application:"
echo "  1. Start MongoDB (if not running)"
echo "  2. Run: make run"
echo ""
echo "Available endpoints:"
echo "  - gRPC: localhost:9090"
echo "  - REST API: http://localhost:8080/api/v1"
echo "  - Health check: http://localhost:8080/health"
echo "  - OpenAPI docs: ./docs/user.swagger.json"
echo ""
echo "Buf CLI commands:"
echo "  - Generate: make buf-generate"
echo "  - Lint: make buf-lint"
echo "  - Format: make buf-format"
echo ""
echo "Example API calls:"
echo "  - Create user: curl -X POST http://localhost:8080/api/v1/users -d '{\"name\":\"John Doe\",\"email\":\"john@example.com\",\"phone\":\"+1234567890\"}'"
echo "  - List users: curl http://localhost:8080/api/v1/users"
