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

	// SQLite specific
	SQLiteFile string `mapstructure:"sqlite_file"`

	// Oracle specific
	ServiceName string `mapstructure:"service_name"`
	SID         string `mapstructure:"sid"`

	// Connection pool settings
	MaxOpenConns    int `mapstructure:"max_open_conns"`
	MaxIdleConns    int `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int `mapstructure:"conn_max_lifetime"`  // seconds
	ConnMaxIdleTime int `mapstructure:"conn_max_idle_time"` // seconds
}

// LoggerConfig represents logger configuration
type LoggerConfig struct {
	Level       string   `mapstructure:"level"`
	OutputPaths []string `mapstructure:"output_paths"`
}
