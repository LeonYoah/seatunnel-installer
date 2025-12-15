package repository

import (
	"context"
	"errors"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"gorm.io/gorm"
)

// clusterRepository 集群Repository实现
type clusterRepository struct {
	*baseRepository[models.Cluster]
}

// newClusterRepository 创建集群Repository
func newClusterRepository(db *gorm.DB) ClusterRepository {
	return &clusterRepository{
		baseRepository: newBaseRepository[models.Cluster](db),
	}
}

// GetByWorkspaceID 根据工作空间ID获取集群列表
func (r *clusterRepository) GetByWorkspaceID(ctx context.Context, workspaceID string, offset, limit int) ([]*models.Cluster, int64, error) {
	var clusters []*models.Cluster
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Cluster{}).Where("workspace_id = ?", workspaceID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Offset(offset).Limit(limit).Find(&clusters).Error
	if err != nil {
		return nil, 0, err
	}

	return clusters, total, nil
}

// GetByName 根据工作空间ID和名称获取集群
func (r *clusterRepository) GetByName(ctx context.Context, workspaceID, name string) (*models.Cluster, error) {
	var cluster models.Cluster
	err := r.db.WithContext(ctx).Where("workspace_id = ? AND name = ?", workspaceID, name).First(&cluster).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cluster, nil
}

// UpdateStatus 更新集群状态
func (r *clusterRepository) UpdateStatus(ctx context.Context, id, status string) error {
	return r.db.WithContext(ctx).Model(&models.Cluster{}).Where("id = ?", id).Update("status", status).Error
}
