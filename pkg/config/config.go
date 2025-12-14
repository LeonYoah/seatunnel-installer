package config

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logger   LoggerConfig   `mapstructure:"logger"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Type     string `mapstructure:"type"` // sqlite, mysql, postgres, oracle
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

// LoggerConfig represents logger configuration
type LoggerConfig struct {
	Level       string   `mapstructure:"level"`
	OutputPaths []string `mapstructure:"output_paths"`
}
