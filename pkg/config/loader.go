package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set config file
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("/etc/seatunnel/")
		v.AddConfigPath("$HOME/.seatunnel/")
	}

	// Read environment variables
	v.SetEnvPrefix("SEATUNNEL")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Set defaults
	setDefaults(v)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found; use defaults
	}

	// Unmarshal config
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate config
	if err := validate(&config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "0.0.0.0")

	// Database defaults
	v.SetDefault("database.type", "sqlite")
	v.SetDefault("database.sqlite_file", "data/seatunnel.db")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", 3600)
	v.SetDefault("database.conn_max_idle_time", 1800)

	// Logger defaults
	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.output_paths", []string{"stdout"})

	// JWT defaults
	v.SetDefault("jwt.secret_key", "seatunnel-enterprise-platform-jwt-secret-key-change-in-production")
	v.SetDefault("jwt.access_token_ttl", 60)   // 60 minutes
	v.SetDefault("jwt.refresh_token_ttl", 168) // 168 hours (7 days)
}

// validate validates the configuration
func validate(config *Config) error {
	// Validate server port
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	// Validate database type
	validDBTypes := map[string]bool{
		"sqlite":   true,
		"mysql":    true,
		"postgres": true,
		"oracle":   true,
	}
	if !validDBTypes[config.Database.Type] {
		return fmt.Errorf("invalid database type: %s", config.Database.Type)
	}

	// Validate logger level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}
	if !validLogLevels[config.Logger.Level] {
		return fmt.Errorf("invalid logger level: %s", config.Logger.Level)
	}

	return nil
}
