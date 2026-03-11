package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration values for the Digna application
type Config struct {
	// Server configuration
	Port string

	// Data directory configuration
	DataDir string

	// Logging configuration
	LogLevel string

	// Database configuration
	SQLiteDriver string
}

// Default configuration values
const (
	DefaultPort         = "8090"
	DefaultDataDir      = "./data/entities"
	DefaultLogLevel     = "info"
	DefaultSQLiteDriver = "sqlite3"
)

// Environment variable names
const (
	EnvPort     = "DIGNA_PORT"
	EnvDataDir  = "DIGNA_DATA_DIR"
	EnvLogLevel = "DIGNA_LOG_LEVEL"
)

// Load loads configuration from environment variables with fallback to defaults
func Load() *Config {
	cfg := &Config{
		Port:         getEnv(EnvPort, DefaultPort),
		DataDir:      getEnv(EnvDataDir, DefaultDataDir),
		LogLevel:     getEnv(EnvLogLevel, DefaultLogLevel),
		SQLiteDriver: DefaultSQLiteDriver,
	}

	// Ensure DataDir has proper format (no trailing slash)
	cfg.DataDir = strings.TrimSuffix(cfg.DataDir, "/")

	return cfg
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetPort returns the port as integer
func (c *Config) GetPort() int {
	port, err := strconv.Atoi(c.Port)
	if err != nil {
		// If port is not a valid number, use default
		defaultPort, _ := strconv.Atoi(DefaultPort)
		return defaultPort
	}
	return port
}

// GetPortString returns the port as string (with colon prefix for http server)
func (c *Config) GetPortString() string {
	return ":" + c.Port
}
