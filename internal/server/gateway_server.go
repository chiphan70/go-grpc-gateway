package server

import (
	"context"
	"fmt"
	"net/http"

	"go-grpcgateway/internal/config"
	"go-grpcgateway/pkg/pb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GatewayServer wraps the HTTP gateway server
type GatewayServer struct {
	server *http.Server
	config *config.Config
}

// NewGatewayServer creates a new HTTP gateway server
func NewGatewayServer(cfg *config.Config) *GatewayServer {
	return &GatewayServer{
		config: cfg,
	}
}

// Start starts the HTTP gateway server
func (s *GatewayServer) Start() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create a connection to the gRPC server
	grpcAddress := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.GRPCPort)
	conn, err := grpc.DialContext(
		ctx,
		grpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("failed to dial gRPC server: %v", err)
	}
	defer conn.Close()

	// Create gRPC-Gateway mux
	mux := runtime.NewServeMux(
		runtime.WithHealthEndpointAt("/health"),
	)

	// Register services
	err = pb.RegisterUserServiceHandler(ctx, mux, conn)
	if err != nil {
		return fmt.Errorf("failed to register user service handler: %v", err)
	}

	// Add CORS middleware
	handler := s.corsMiddleware(mux)

	// Create HTTP server
	address := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.HTTPPort)
	s.server = &http.Server{
		Addr:    address,
		Handler: handler,
	}

	logrus.WithField("address", address).Info("Starting HTTP gateway server")

	return s.server.ListenAndServe()
}

// Stop stops the HTTP gateway server gracefully
func (s *GatewayServer) Stop() error {
	logrus.Info("Stopping HTTP gateway server")
	if s.server != nil {
		return s.server.Shutdown(context.Background())
	}
	return nil
}

// corsMiddleware adds CORS headers to the response
func (s *GatewayServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
