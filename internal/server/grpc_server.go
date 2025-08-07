package server

import (
	"fmt"
	"net"

	"go-grpcgateway/internal/config"
	"go-grpcgateway/internal/db"
	"go-grpcgateway/internal/service"
	"go-grpcgateway/pkg/pb"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCServer wraps the gRPC server
type GRPCServer struct {
	server      *grpc.Server
	userService *service.UserService
	config      *config.Config
}

// NewGRPCServer creates a new gRPC server
func NewGRPCServer(cfg *config.Config, database *db.MongoDB) *GRPCServer {
	server := grpc.NewServer()

	// Enable reflection for development
	reflection.Register(server)

	// Create services
	userService := service.NewUserService(database)

	// Register services
	pb.RegisterUserServiceServer(server, userService)

	return &GRPCServer{
		server:      server,
		userService: userService,
		config:      cfg,
	}
}

// Start starts the gRPC server
func (s *GRPCServer) Start() error {
	address := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.GRPCPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	logrus.WithField("address", address).Info("Starting gRPC server")

	return s.server.Serve(lis)
}

// Stop stops the gRPC server gracefully
func (s *GRPCServer) Stop() {
	logrus.Info("Stopping gRPC server")
	s.server.GracefulStop()
}

// GetServer returns the underlying gRPC server
func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}
