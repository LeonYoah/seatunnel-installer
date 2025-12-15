package repository

import (
	"context"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
)

// BaseRepository 基础Repository接口
type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id string) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]*T, int64, error)
}

// TenantRepository 租户Repository接口
type TenantRepository interface {
	BaseRepository[models.Tenant]
	GetByName(ctx context.Context, name string) (*models.Tenant, error)
	ListActive(ctx context.Context, offset, limit int) ([]*models.Tenant, int64, error)
}

// WorkspaceRepository 工作空间Repository接口
type WorkspaceRepository interface {
	BaseRepository[models.Workspace]
	GetByTenantID(ctx context.Context, tenantID string, offset, limit int) ([]*models.Workspace, int64, error)
	GetByTenantAndName(ctx context.Context, tenantID, name string) (*models.Workspace, error)
}

// HostRepository 主机Repository接口
type HostRepository interface {
	BaseRepository[models.Host]
	GetByWorkspaceID(ctx context.Context, workspaceID string, offset, limit int) ([]*models.Host, int64, error)
	GetByIP(ctx context.Context, workspaceID, ip string) (*models.Host, error)
	UpdateAgentStatus(ctx context.Context, id, status string) error
	UpdateHeartbeat(ctx context.Context, id string) error
	GetOnlineHosts(ctx context.Context, workspaceID string) ([]*models.Host, error)
}

// ClusterRepository 集群Repository接口
type ClusterRepository interface {
	BaseRepository[models.Cluster]
	GetByWorkspaceID(ctx context.Context, workspaceID string, offset, limit int) ([]*models.Cluster, int64, error)
	GetByName(ctx context.Context, workspaceID, name string) (*models.Cluster, error)
	UpdateStatus(ctx context.Context, id, status string) error
}

// NodeRepository 节点Repository接口
type NodeRepository interface {
	BaseRepository[models.Node]
	GetByClusterID(ctx context.Context, clusterID string) ([]*models.Node, error)
	GetByIP(ctx context.Context, clusterID, ip string) (*models.Node, error)
	UpdateStatus(ctx context.Context, id, status string) error
	UpdateHeartbeat(ctx context.Context, id string) error
}

// TaskRepository 任务Repository接口
type TaskRepository interface {
	BaseRepository[models.Task]
	GetByWorkspaceID(ctx context.Context, workspaceID string, offset, limit int) ([]*models.Task, int64, error)
	GetByName(ctx context.Context, workspaceID, name string) (*models.Task, error)
	GetVersions(ctx context.Context, taskID string) ([]*models.Task, error)
	GetLatestVersion(ctx context.Context, taskID string) (*models.Task, error)
}

// RunRepository 运行Repository接口
type RunRepository interface {
	BaseRepository[models.Run]
	GetByTaskID(ctx context.Context, taskID string, offset, limit int) ([]*models.Run, int64, error)
	GetByClusterID(ctx context.Context, clusterID string, offset, limit int) ([]*models.Run, int64, error)
	UpdateStatus(ctx context.Context, id, status string) error
	GetRunningRuns(ctx context.Context) ([]*models.Run, error)
}

// AuditLogRepository 审计日志Repository接口
type AuditLogRepository interface {
	Create(ctx context.Context, log *models.AuditLog) error
	GetByTenantID(ctx context.Context, tenantID string, offset, limit int) ([]*models.AuditLog, int64, error)
	GetByUserID(ctx context.Context, userID string, offset, limit int) ([]*models.AuditLog, int64, error)
	GetByResource(ctx context.Context, tenantID, resource, resourceID string, offset, limit int) ([]*models.AuditLog, int64, error)
}

// SecretRepository 凭证Repository接口
type SecretRepository interface {
	BaseRepository[models.Secret]
	GetByWorkspaceID(ctx context.Context, workspaceID string, offset, limit int) ([]*models.Secret, int64, error)
	GetByName(ctx context.Context, workspaceID, name string) (*models.Secret, error)
}

// TransactionManager 事务管理器接口
type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// RepositoryManager Repository管理器接口
type RepositoryManager interface {
	TransactionManager
	Tenant() TenantRepository
	Workspace() WorkspaceRepository
	Host() HostRepository
	Cluster() ClusterRepository
	Node() NodeRepository
	Task() TaskRepository
	Run() RunRepository
	AuditLog() AuditLogRepository
	Secret() SecretRepository
}
