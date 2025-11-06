package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds application configuration
type Config struct {
	Server ServerConfig
	Client ClientConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port                string
	MaxConcurrentStreams uint32
	MaxMessageSize       int
}

// ClientConfig holds client-specific configuration
type ClientConfig struct {
	ServerAddress    string
	ConnectionTimeout time.Duration
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:                getEnv("GRPC_PORT", ":50051"),
			MaxConcurrentStreams: getEnvAsUint32("MAX_CONCURRENT_STREAMS", 1000),
			MaxMessageSize:       getEnvAsInt("MAX_MESSAGE_SIZE", 4*1024*1024), // 4MB
		},
		Client: ClientConfig{
			ServerAddress:    getEnv("GRPC_SERVER_ADDRESS", "localhost:50051"),
			ConnectionTimeout: getEnvAsDuration("CONNECTION_TIMEOUT", 5*time.Second),
		},
	}
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsUint32(key string, defaultValue uint32) uint32 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseUint(value, 10, 32); err == nil {
			return uint32(intValue)
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}