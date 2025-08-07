package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go-grpcgateway/internal/config"
	"go-grpcgateway/internal/db"
	"go-grpcgateway/internal/server"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "go-grpcgateway",
		Short: "A gRPC Gateway service with MongoDB",
		Long:  "A complete gRPC service with REST API gateway and MongoDB integration",
		Run:   runServer,
	}

	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("Failed to execute command")
	}
}

func runServer(cmd *cobra.Command, args []string) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration")
	}

	// Connect to MongoDB
	mongodb, err := db.NewMongoDB(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to MongoDB")
	}
	defer func() {
		if err := mongodb.Close(); err != nil {
			logrus.WithError(err).Error("Failed to close MongoDB connection")
		}
	}()

	// Create servers
	grpcServer := server.NewGRPCServer(cfg, mongodb)
	gatewayServer := server.NewGatewayServer(cfg)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start servers in goroutines
	var wg sync.WaitGroup

	// Start gRPC server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := grpcServer.Start(); err != nil {
			logrus.WithError(err).Error("gRPC server failed")
		}
	}()

	// Start HTTP gateway server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := gatewayServer.Start(); err != nil {
			logrus.WithError(err).Error("HTTP gateway server failed")
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		logrus.WithField("signal", sig).Info("Received shutdown signal")
	case <-ctx.Done():
		logrus.Info("Context cancelled")
	}

	// Graceful shutdown
	logrus.Info("Shutting down servers...")

	// Stop HTTP gateway server
	if err := gatewayServer.Stop(); err != nil {
		logrus.WithError(err).Error("Failed to stop HTTP gateway server")
	}

	// Stop gRPC server
	grpcServer.Stop()

	// Wait for all goroutines to finish
	wg.Wait()

	logrus.Info("Servers stopped successfully")
}
