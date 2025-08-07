package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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
	mux := runtime.NewServeMux()

	// Register services
	err = pb.RegisterUserServiceHandler(ctx, mux, conn)
	if err != nil {
		return fmt.Errorf("failed to register user service handler: %v", err)
	}

	// Add Swagger UI and CORS middleware
	handler := s.setupRoutes(mux)

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

// setupRoutes configures all routes including Swagger UI
func (s *GatewayServer) setupRoutes(mux *runtime.ServeMux) http.Handler {
	// Create main router
	mainMux := http.NewServeMux()

	// Serve gRPC-Gateway
	mainMux.Handle("/", s.corsMiddleware(mux))

	// Serve Swagger JSON
	mainMux.HandleFunc("/swagger.json", s.swaggerJSONHandler)

	// Serve Swagger UI
	mainMux.HandleFunc("/swagger/", s.swaggerUIHandler)

	return mainMux
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

// swaggerJSONHandler serves the OpenAPI JSON file
func (s *GatewayServer) swaggerJSONHandler(w http.ResponseWriter, r *http.Request) {
	swaggerFile := filepath.Join("docs", "api.swagger.json")
	data, err := os.ReadFile(swaggerFile)
	if err != nil {
		http.Error(w, "Swagger file not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// swaggerUIHandler serves the Swagger UI
func (s *GatewayServer) swaggerUIHandler(w http.ResponseWriter, r *http.Request) {
	swaggerHTML := `
	<!DOCTYPE html>
		<html>
		<head>
			<title>API Documentation</title>
			<link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui.css" />
			<style>
				html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
				*, *:before, *:after { box-sizing: inherit; }
				body { margin:0; background: #fafafa; }
			</style>
		</head>
		<body>
			<div id="swagger-ui"></div>
			<script src="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-bundle.js"></script>
			<script src="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-standalone-preset.js"></script>
			<script>
				window.onload = function() {
					const ui = SwaggerUIBundle({
						url: '/swagger.json',
						dom_id: '#swagger-ui',
						deepLinking: true,
						presets: [
							SwaggerUIBundle.presets.apis,
							SwaggerUIStandalonePreset
						],
						plugins: [
							SwaggerUIBundle.plugins.DownloadUrl
						],
						layout: "StandaloneLayout"
					});
				}
			</script>
		</body>
	</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(swaggerHTML))
}
