package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config holds all configuration for our application
type Config struct {
	MongoDB  MongoDBConfig
	Server   ServerConfig
	LogLevel string
}

// MongoDBConfig holds MongoDB configuration
type MongoDBConfig struct {
	URI      string
	Database string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	GRPCPort string
	HTTPPort string
	Host     string
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found")
	}

	// Set default values
	viper.SetDefault("MONGODB_URI", "mongodb://localhost:27017")
	viper.SetDefault("MONGODB_DATABASE", "grpcgateway_db")
	viper.SetDefault("GRPC_PORT", "9090")
	viper.SetDefault("HTTP_PORT", "8080")
	viper.SetDefault("HOST", "localhost")
	viper.SetDefault("LOG_LEVEL", "info")

	// Bind environment variables
	viper.AutomaticEnv()

	// Get values from environment or defaults
	config := &Config{
		MongoDB: MongoDBConfig{
			URI:      viper.GetString("MONGODB_URI"),
			Database: viper.GetString("MONGODB_DATABASE"),
		},
		Server: ServerConfig{
			GRPCPort: viper.GetString("GRPC_PORT"),
			HTTPPort: viper.GetString("HTTP_PORT"),
			Host:     viper.GetString("HOST"),
		},
		LogLevel: viper.GetString("LOG_LEVEL"),
	}

	// Set log level
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logrus.Warn("Invalid log level, using info")
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	return config, nil
}

// GetEnv gets an environment variable with a fallback value
func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
