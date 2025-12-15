package repository

import (
	"context"
	"errors"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"gorm.io/gorm"
)

// workspaceRepository 工作空间Repository实现
type workspaceRepository struct {
	*baseRepository[models.Workspace]
}

// newWorkspaceRepository 创建工作空间Repository
func newWorkspaceRepository(db *gorm.DB) WorkspaceRepository {
	return &workspaceRepository{
		baseRepository: newBaseRepository[models.Workspace](db),
	}
}

// GetByTenantID 根据租户ID获取工作空间列表
func (r *workspaceRepository) GetByTenantID(ctx context.Context, tenantID string, offset, limit int) ([]*models.Workspace, int64, error) {
	var workspaces []*models.Workspace
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Workspace{}).Where("tenant_id = ?", tenantID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Offset(offset).Limit(limit).Find(&workspaces).Error
	if err != nil {
		return nil, 0, err
	}

	return workspaces, total, nil
}

// GetByTenantAndName 根据租户ID和名称获取工作空间
func (r *workspaceRepository) GetByTenantAndName(ctx context.Context, tenantID, name string) (*models.Workspace, error) {
	var workspace models.Workspace
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND name = ?", tenantID, name).First(&workspace).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &workspace, nil
}
