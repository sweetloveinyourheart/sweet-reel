package config

import (
	"time"
)

const (
	// Server defaults
	DefaultHost            = "0.0.0.0"
	DefaultReadTimeout     = 30 * time.Second
	DefaultWriteTimeout    = 30 * time.Second
	DefaultShutdownTimeout = 10 * time.Second
	DefaultBodyLimit       = 4 * 1024 * 1024

	// Security defaults
	DefaultTokenExpiration = 3600 * time.Second
	DefaultAllowOrigin     = "*"

	// Logging defaults
	DefaultLogLevel   = "info"
	DefaultLogFormat  = "json"
	DefaultRequestLog = true
)

// Config holds all configuration for the API Gateway service
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Security SecurityConfig `mapstructure:"security"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port            uint64        `mapstructure:"port"`
	Host            string        `mapstructure:"host"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	BodyLimit       int           `mapstructure:"body_limit"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	JWTSecret       string        `mapstructure:"jwt_secret"`
	TokenExpiration time.Duration `mapstructure:"token_expiration"`
	AllowOrigins    []string      `mapstructure:"allow_origins"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	RequestLog bool   `mapstructure:"request_log"`
}

func LoadServerConfig(port uint64, signingKey string) Config {
	return Config{
		Server: ServerConfig{
			Port:            port,
			Host:            DefaultHost,
			ReadTimeout:     DefaultReadTimeout,
			WriteTimeout:    DefaultWriteTimeout,
			ShutdownTimeout: DefaultShutdownTimeout,
			BodyLimit:       DefaultBodyLimit,
		},
		Security: SecurityConfig{
			JWTSecret:       signingKey,
			TokenExpiration: DefaultTokenExpiration,
			AllowOrigins:    []string{DefaultAllowOrigin},
		},
		Logging: LoggingConfig{
			Level:      DefaultLogLevel,
			Format:     DefaultLogFormat,
			RequestLog: DefaultRequestLog,
		},
	}

}
