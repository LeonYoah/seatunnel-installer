package config

import (
	"time"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/database"
)

// ToDatabaseConfig converts application config to database config
func (c *Config) ToDatabaseConfig() *database.Config {
	return &database.Config{
		Type:            c.Database.Type,
		Host:            c.Database.Host,
		Port:            c.Database.Port,
		Database:        c.Database.Database,
		Username:        c.Database.Username,
		Password:        c.Database.Password,
		SQLiteFile:      c.Database.SQLiteFile,
		ServiceName:     c.Database.ServiceName,
		SID:             c.Database.SID,
		SSLMode:         c.Database.SSLMode,
		MaxOpenConns:    c.Database.MaxOpenConns,
		MaxIdleConns:    c.Database.MaxIdleConns,
		ConnMaxLifetime: time.Duration(c.Database.ConnMaxLifetime) * time.Second,
		ConnMaxIdleTime: time.Duration(c.Database.ConnMaxIdleTime) * time.Second,
	}
}
