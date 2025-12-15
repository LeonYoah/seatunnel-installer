package repository

import (
	"context"

	"gorm.io/gorm"
)

// contextKey 定义上下文键类型，避免字符串键冲突
type contextKey string

const (
	txManagerKey contextKey = "tx_manager"
)

// repositoryManager Repository管理器实现
type repositoryManager struct {
	db *gorm.DB

	// Repository实例
	tenantRepo    TenantRepository
	workspaceRepo WorkspaceRepository
	hostRepo      HostRepository
	clusterRepo   ClusterRepository
	nodeRepo      NodeRepository
	taskRepo      TaskRepository
	runRepo       RunRepository
	auditLogRepo  AuditLogRepository
	secretRepo    SecretRepository
}

// NewRepositoryManager 创建Repository管理器
func NewRepositoryManager(db *gorm.DB) RepositoryManager {
	return &repositoryManager{
		db:            db,
		tenantRepo:    newTenantRepository(db),
		workspaceRepo: newWorkspaceRepository(db),
		hostRepo:      newHostRepository(db),
		clusterRepo:   newClusterRepository(db),
		nodeRepo:      newNodeRepository(db),
		taskRepo:      newTaskRepository(db),
		runRepo:       newRunRepository(db),
		auditLogRepo:  newAuditLogRepository(db),
		secretRepo:    newSecretRepository(db),
	}
}

// WithTransaction 执行事务
func (m *repositoryManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建新的Repository管理器，使用事务连接
		txManager := &repositoryManager{
			db:            tx,
			tenantRepo:    newTenantRepository(tx),
			workspaceRepo: newWorkspaceRepository(tx),
			hostRepo:      newHostRepository(tx),
			clusterRepo:   newClusterRepository(tx),
			nodeRepo:      newNodeRepository(tx),
			taskRepo:      newTaskRepository(tx),
			runRepo:       newRunRepository(tx),
			auditLogRepo:  newAuditLogRepository(tx),
			secretRepo:    newSecretRepository(tx),
		}

		// 将事务管理器放入上下文
		txCtx := context.WithValue(ctx, txManagerKey, txManager)
		return fn(txCtx)
	})
}

// Tenant 获取租户Repository
func (m *repositoryManager) Tenant() TenantRepository {
	return m.tenantRepo
}

// Workspace 获取工作空间Repository
func (m *repositoryManager) Workspace() WorkspaceRepository {
	return m.workspaceRepo
}

// Host 获取主机Repository
func (m *repositoryManager) Host() HostRepository {
	return m.hostRepo
}

// Cluster 获取集群Repository
func (m *repositoryManager) Cluster() ClusterRepository {
	return m.clusterRepo
}

// Node 获取节点Repository
func (m *repositoryManager) Node() NodeRepository {
	return m.nodeRepo
}

// Task 获取任务Repository
func (m *repositoryManager) Task() TaskRepository {
	return m.taskRepo
}

// Run 获取运行Repository
func (m *repositoryManager) Run() RunRepository {
	return m.runRepo
}

// AuditLog 获取审计日志Repository
func (m *repositoryManager) AuditLog() AuditLogRepository {
	return m.auditLogRepo
}

// Secret 获取凭证Repository
func (m *repositoryManager) Secret() SecretRepository {
	return m.secretRepo
}

// GetRepositoryManager 从上下文获取Repository管理器
func GetRepositoryManager(ctx context.Context) RepositoryManager {
	if txManager, ok := ctx.Value(txManagerKey).(RepositoryManager); ok {
		return txManager
	}
	// 如果上下文中没有事务管理器，返回nil或抛出错误
	// 这里简单返回nil，实际使用时应该从依赖注入容器获取
	return nil
}
