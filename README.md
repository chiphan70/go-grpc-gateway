# Go gRPC-Gateway with MongoDB

A complete Golang template project that demonstrates how to build a gRPC service with REST API gateway and MongoDB integration. This project automatically generates REST APIs from gRPC service definitions using grpc-gateway.

## Features

- 🚀 **gRPC Server**: High-performance gRPC service
- 🌐 **REST API Gateway**: Automatic REST API generation from gRPC definitions
- 🗄️ **MongoDB Integration**: Complete CRUD operations with MongoDB
- 📝 **Protocol Buffers**: Type-safe API definitions with Buf CLI
- 🐳 **Docker Support**: Containerized deployment with Docker Compose
- ⚡ **Hot Reload**: Development environment with live reload
- 📊 **Health Checks**: Built-in health check endpoints
- 🔧 **Configuration**: Environment-based configuration
- 📚 **OpenAPI/Swagger**: Auto-generated API documentation
- 🧪 **Testing**: Test structure and examples
- 🛠️ **Buf CLI**: Modern protobuf toolchain for linting, formatting, and breaking change detection

## Project Structure (Go Standard Layout)

```
go-grpcgateway/
├── api/
│   └── proto/              # Protocol Buffer definitions (.proto files)
├── cmd/                    # Application entrypoints (main applications)
├── internal/               # Private application code (internal packages)
│   ├── config/            # Configuration management
│   ├── db/                # Database connection and utilities
│   ├── models/            # Data models and domain entities
│   ├── repository/        # Data access layer (repositories)
│   ├── server/            # Server implementations (gRPC, HTTP)
│   └── service/           # Business logic and service layer
├── pkg/                   # Public packages (can be imported by external apps)
│   └── pb/               # Generated protobuf files
├── scripts/              # Build, setup, and deployment scripts
├── docker-compose.yml    # Docker Compose configuration
├── Dockerfile           # Docker build configuration
├── Makefile            # Build automation
└── README.md           # This file
```

### **About Directory Structure:**

#### **`internal/` - Private Packages**

- **Purpose**: Contains code for internal use within this project only
- **Characteristic**: Go compiler prevents external packages from importing code from `internal/`
- **Contains**: Business logic, implementation details, internal services
- **Examples**:
  - `internal/config` - Configuration management
  - `internal/service` - Business logic layer
  - `internal/models` - Domain models
  - `internal/db` - Database layer

#### **`pkg/` - Public Packages**

- **Purpose**: Contains code that can be reused by external applications
- **Characteristic**: Public API, stable interfaces
- **Contains**: Libraries, utilities, generated code
- **Examples**:
  - `pkg/pb` - Generated protobuf files (can be imported by clients)

#### **Why This Structure is Better:**

1. **Security**: `internal/` prevents external packages from accessing implementation details
2. **Maintainability**: Clear separation between public API and private implementation
3. **Modularity**: Easy to refactor internal code without affecting external dependencies
4. **Go Conventions**: Follows standard Go project layout conventions

## Prerequisites

- Go 1.21 or higher
- **Buf CLI** (recommended) or Protocol Buffers compiler (protoc)
- MongoDB (local or Docker)
- Docker and Docker Compose (optional)

### Installing Prerequisites

#### macOS

```bash
# Install Go
brew install go

# Install Buf CLI (recommended)
brew install bufbuild/buf/buf

# Install protoc (optional, for legacy support)
brew install protobuf

# Install MongoDB
brew tap mongodb/brew
brew install mongodb-community
```

#### Ubuntu/Debian

```bash
# Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install Buf CLI (recommended)
curl -sSL "https://github.com/bufbuild/buf/releases/latest/download/buf-$(uname -s)-$(uname -m)" -o "/tmp/buf"
chmod +x "/tmp/buf"
sudo mv "/tmp/buf" "/usr/local/bin/buf"

# Install protoc (optional, for legacy support)
sudo apt-get update
sudo apt-get install protobuf-compiler

# Install MongoDB
wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | sudo apt-key add -
echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/6.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-6.0.list
sudo apt-get update
sudo apt-get install mongodb-org
```

## Quick Start

### 1. Clone and Setup

```bash
# Clone the repository
git clone <your-repo-url>
cd go-grpcgateway

# Run the setup script
./scripts/setup.sh
```

### 2. Configure Environment

Copy the environment file and adjust settings if needed:

```bash
cp .env.example .env
```

Edit `.env` file:

```env
# Database Configuration
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=grpcgateway_db

# Server Configuration
GRPC_PORT=9090
HTTP_PORT=8080
HOST=localhost

# Log Configuration
LOG_LEVEL=info
```

### 3. Start MongoDB

```bash
# Using Homebrew (macOS)
brew services start mongodb/brew/mongodb-community

# Using systemd (Linux)
sudo systemctl start mongod

# Using Docker
docker run -d -p 27017:27017 --name mongodb mongo:6.0
```

### 4. Run the Application

```bash
# Build and run
make run

# Or run with live reload (requires air)
make dev
```

## API Endpoints

Once the application is running, you can access:

- **gRPC Server**: `localhost:9090`
- **REST API**: `http://localhost:8080/api/v1`
- **Swagger UI**: `http://localhost:8080/swagger/` (Interactive API Documentation)
- **OpenAPI JSON**: `http://localhost:8080/swagger.json` (Raw OpenAPI specification)

### REST API Examples

#### Create a User

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1234567890"
  }'
```

#### Get All Users

```bash
curl http://localhost:8080/api/v1/users
```

#### Get User by ID

```bash
curl http://localhost:8080/api/v1/users/{user_id}
```

#### Update User

```bash
curl -X PUT http://localhost:8080/api/v1/users/{user_id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "email": "john.smith@example.com",
    "phone": "+1234567890"
  }'
```

#### Delete User

```bash
curl -X DELETE http://localhost:8080/api/v1/users/{user_id}
```

### gRPC Examples

You can test gRPC endpoints using tools like [grpcurl](https://github.com/fullstorydev/grpcurl):

```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# List services
grpcurl -plaintext localhost:9090 list

# Create a user
grpcurl -plaintext -d '{"name":"John Doe","email":"john@example.com","phone":"+1234567890"}' \
  localhost:9090 user.v1.UserService/CreateUser

# List users
grpcurl -plaintext localhost:9090 user.v1.UserService/ListUsers
```

## Development

### Available Make Commands

#### Basic Commands

```bash
make help           # Show all available commands
make deps           # Download dependencies
make build          # Build the application
make run            # Build and run the application
make dev            # Run with live reload
make test           # Run tests
make test-coverage  # Run tests with coverage
make clean          # Clean build artifacts
```

#### Buf CLI Commands (Recommended)

```bash
make install-buf    # Install Buf CLI
make buf-generate   # Generate protobuf files with Buf
make buf-lint       # Lint protobuf files
make buf-format     # Format protobuf files
make buf-breaking   # Check for breaking changes
make buf-deps       # Update Buf dependencies
```

#### Legacy Protoc Commands

```bash
make tools          # Install required protoc tools
make proto          # Generate protobuf files with protoc
```

### Adding New Services

1. **Define the service in Protocol Buffers**:

   Create a new `.proto` file in `api/proto/`:

   ```protobuf
   // api/proto/product.proto
   syntax = "proto3";

   package product.v1;

   option go_package = "go-grpcgateway/pkg/pb";

   import "google/api/annotations.proto";

   service ProductService {
     rpc CreateProduct(CreateProductRequest) returns (Product) {
       option (google.api.http) = {
         post: "/api/v1/products"
         body: "*"
       };
     }
   }
   ```

2. **Create the model**:

   ```go
   // internal/models/product.go
   package models

   type Product struct {
       ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
       Name  string             `bson:"name" json:"name"`
       Price float64            `bson:"price" json:"price"`
   }
   ```

3. **Implement the service**:

   ```go
   // internal/service/product_service.go
   package service

   type ProductService struct {
       pb.UnimplementedProductServiceServer
       // ... implementation
   }
   ```

4. **Register the service** in `internal/server/grpc_server.go`

5. **Generate protobuf files**:
   ```bash
   make buf-generate  # Recommended
   # or
   make proto         # Legacy protoc
   ```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific tests
go test ./internal/service/...
```

## Docker Deployment

### Using Docker Compose (Recommended)

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Using Docker Only

```bash
# Build the image
make docker-build

# Run MongoDB
docker run -d -p 27017:27017 --name mongodb mongo:6.0

# Run the application
make docker-run
```

## Architecture Overview

### **Layer Architecture**

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                       │
│                  (gRPC + REST Gateway)                      │
├─────────────────────────────────────────────────────────────┤
│                     Service Layer                          │
│                  (Business Logic)                          │
├─────────────────────────────────────────────────────────────┤
│                   Repository Layer                         │
│                   (Data Access)                            │
├─────────────────────────────────────────────────────────────┤
│                    Database Layer                          │
│                     (MongoDB)                              │
└─────────────────────────────────────────────────────────────┘
```

### **Request Flow**

```
HTTP Request → gRPC Gateway → gRPC Server → Service → Repository → MongoDB
            ←               ←             ←         ←            ←
```

## Production Considerations

### Security

- Use authentication middleware
- Implement rate limiting
- Use HTTPS/TLS certificates
- Validate all inputs
- Use environment variables for sensitive data

### Performance

- Implement connection pooling
- Add caching layer (Redis)
- Use MongoDB indexes
- Implement graceful shutdown
- Monitor with metrics

### Monitoring

- Add structured logging
- Implement health checks
- Use metrics collection (Prometheus)
- Add distributed tracing

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Buf CLI Features

### Why Buf CLI?

Buf CLI provides several advantages over traditional protoc:

1. **🚀 Performance**: Faster compilation and parallel processing
2. **🔍 Better Linting**: More comprehensive linting rules
3. **📦 Dependency Management**: Simplified dependency management with buf.build registry
4. **🔄 Breaking Change Detection**: Automated API compatibility checking
5. **📝 Better Error Messages**: More descriptive and actionable error messages
6. **🌐 Remote Generation**: Generate code without installing plugins locally

### Buf Configuration Files

- **`buf.yaml`**: Main configuration file (linting rules, dependencies)
- **`buf.gen.yaml`**: Code generation configuration (plugins, output paths)
- **`buf.work.yaml`**: Workspace configuration (multi-module projects)

### Buf Commands

```bash
# Generate code
buf generate

# Lint protobuf files
buf lint

# Format protobuf files
buf format -w

# Check for breaking changes
buf breaking --against '.git#branch=main'

# Update dependencies
buf mod update

# Push to buf.build registry (if public)
buf push
```

## Resources

- [Go Standard Project Layout](https://github.com/golang-standards/project-layout)
- [Buf CLI Documentation](https://docs.buf.build/)
- [Buf Schema Registry](https://buf.build/)
- [gRPC](https://grpc.io/)
- [gRPC-Gateway](https://grpc-ecosystem.github.io/grpc-gateway/)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [MongoDB Go Driver](https://pkg.go.dev/go.mongodb.org/mongo-driver)
