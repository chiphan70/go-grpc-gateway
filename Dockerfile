# Build stage
FROM golang:1.21-alpine AS builder

# Install protoc and git
RUN apk add --no-cache protobuf git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install protoc plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# Copy source code
COPY . .

# Generate protobuf files
RUN make proto

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/go-grpcgateway ./cmd

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/bin/go-grpcgateway .

# Expose ports
EXPOSE 8080 9090

# Run the application
CMD ["./go-grpcgateway"]
