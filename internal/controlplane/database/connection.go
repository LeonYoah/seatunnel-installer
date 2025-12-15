package database

import (
	"fmt"
	"os"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connection 数据库连接管理器
type Connection struct {
	db     *gorm.DB
	config *Config
}

// NewConnection 创建数据库连接
func NewConnection(config *Config) (*Connection, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("数据库配置验证失败: %w", err)
	}

	var dialector gorm.Dialector

	switch config.Type {
	case "sqlite":
		var err error
		dialector, err = createSQLiteDialector(config)
		if err != nil {
			return nil, err
		}
	case "mysql":
		dialector = mysql.Open(config.DSN())
	case "postgres":
		dialector = postgres.Open(config.DSN())
	case "oracle":
		var err error
		dialector, err = createOracleDialector(config)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", config.Type)
	}

	// GORM配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层数据库连接失败: %w", err)
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	return &Connection{
		db:     db,
		config: config,
	}, nil
}

// DB 获取GORM数据库实例
func (c *Connection) DB() *gorm.DB {
	return c.db
}

// Config 获取数据库配置
func (c *Connection) Config() *Config {
	return c.config
}

// Migrate 执行数据库迁移
func (c *Connection) Migrate() error {
	return models.AutoMigrate(c.db)
}

// Close 关闭数据库连接
func (c *Connection) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping 测试数据库连接
func (c *Connection) Ping() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// ensureDir 确保目录存在
func ensureDir(dir string) error {
	if dir == "" || dir == "." {
		return nil
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
