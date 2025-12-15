package database

import (
	"fmt"
	"time"
)

// Config 数据库配置
type Config struct {
	Type     string `mapstructure:"type" json:"type"` // sqlite/mysql/postgres/oracle
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Database string `mapstructure:"database" json:"database"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`

	// 连接池配置
	MaxOpenConns    int           `mapstructure:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" json:"conn_max_idle_time"`

	// SQLite特定配置
	SQLiteFile string `mapstructure:"sqlite_file" json:"sqlite_file"`

	// SSL配置
	SSLMode string `mapstructure:"ssl_mode" json:"ssl_mode"`

	// Oracle特定配置
	ServiceName string `mapstructure:"service_name" json:"service_name"` // Oracle服务名
	SID         string `mapstructure:"sid" json:"sid"`                   // Oracle SID
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Type:            "sqlite",
		Host:            "localhost",
		Port:            3306,
		Database:        "seatunnel",
		Username:        "root",
		Password:        "",
		MaxOpenConns:    100,
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: time.Minute * 30,
		SQLiteFile:      "data/seatunnel.db",
		SSLMode:         "disable",
		ServiceName:     "",
		SID:             "",
	}
}

// DefaultOracleConfig 返回Oracle默认配置
func DefaultOracleConfig() *Config {
	return &Config{
		Type:            "oracle",
		Host:            "localhost",
		Port:            1521,
		Database:        "XE", // Oracle Express Edition默认数据库
		Username:        "system",
		Password:        "",
		MaxOpenConns:    50, // Oracle连接数通常较少
		MaxIdleConns:    5,  // Oracle空闲连接数
		ConnMaxLifetime: time.Hour * 2,
		ConnMaxIdleTime: time.Minute * 15,
		ServiceName:     "XE", // Oracle Express Edition默认服务名
		SID:             "",
	}
}

// DSN 生成数据源名称
func (c *Config) DSN() string {
	switch c.Type {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.Username, c.Password, c.Host, c.Port, c.Database)
	case "postgres":
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
			c.Host, c.Username, c.Password, c.Database, c.Port, c.SSLMode)
	case "oracle":
		// Oracle DSN格式: user/password@host:port/service_name 或 user/password@host:port:sid
		if c.ServiceName != "" {
			return fmt.Sprintf("%s/%s@%s:%d/%s",
				c.Username, c.Password, c.Host, c.Port, c.ServiceName)
		} else if c.SID != "" {
			return fmt.Sprintf("%s/%s@%s:%d:%s",
				c.Username, c.Password, c.Host, c.Port, c.SID)
		} else {
			// 默认使用数据库名作为服务名
			return fmt.Sprintf("%s/%s@%s:%d/%s",
				c.Username, c.Password, c.Host, c.Port, c.Database)
		}
	case "sqlite":
		return c.SQLiteFile
	default:
		return ""
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Type == "" {
		return fmt.Errorf("数据库类型不能为空")
	}

	switch c.Type {
	case "sqlite":
		if c.SQLiteFile == "" {
			return fmt.Errorf("SQLite文件路径不能为空")
		}
	case "mysql", "postgres":
		if c.Host == "" {
			return fmt.Errorf("数据库主机不能为空")
		}
		if c.Port <= 0 {
			return fmt.Errorf("数据库端口必须大于0")
		}
		if c.Database == "" {
			return fmt.Errorf("数据库名称不能为空")
		}
		if c.Username == "" {
			return fmt.Errorf("数据库用户名不能为空")
		}
	case "oracle":
		if c.Host == "" {
			return fmt.Errorf("Oracle数据库主机不能为空")
		}
		if c.Port <= 0 {
			c.Port = 1521 // Oracle默认端口
		}
		if c.Username == "" {
			return fmt.Errorf("Oracle数据库用户名不能为空")
		}
		// Oracle必须有ServiceName、SID或Database其中之一
		if c.ServiceName == "" && c.SID == "" && c.Database == "" {
			return fmt.Errorf("Oracle数据库必须指定ServiceName、SID或Database其中之一")
		}
	default:
		return fmt.Errorf("不支持的数据库类型: %s", c.Type)
	}

	return nil
}
