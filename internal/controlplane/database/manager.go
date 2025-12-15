package database

import (
	"fmt"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/repository"
	"gorm.io/gorm"
)

// Manager 数据库管理器
type Manager struct {
	connection *Connection
	repoMgr    repository.RepositoryManager
}

// NewManager 创建数据库管理器
func NewManager(config *Config) (*Manager, error) {
	// 创建数据库连接
	conn, err := NewConnection(config)
	if err != nil {
		return nil, fmt.Errorf("创建数据库连接失败: %w", err)
	}

	// 执行数据库迁移
	if err := conn.Migrate(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("数据库迁移失败: %w", err)
	}

	// 创建Repository管理器
	repoMgr := repository.NewRepositoryManager(conn.DB())

	return &Manager{
		connection: conn,
		repoMgr:    repoMgr,
	}, nil
}

// DB 获取GORM数据库实例
func (m *Manager) DB() *gorm.DB {
	return m.connection.DB()
}

// Repository 获取Repository管理器
func (m *Manager) Repository() repository.RepositoryManager {
	return m.repoMgr
}

// Connection 获取数据库连接
func (m *Manager) Connection() *Connection {
	return m.connection
}

// Close 关闭数据库连接
func (m *Manager) Close() error {
	return m.connection.Close()
}

// Ping 测试数据库连接
func (m *Manager) Ping() error {
	return m.connection.Ping()
}
