package repository

import (
	"context"
	"errors"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"gorm.io/gorm"
)

// secretRepository 凭证Repository实现
type secretRepository struct {
	*baseRepository[models.Secret]
}

// newSecretRepository 创建凭证Repository
func newSecretRepository(db *gorm.DB) SecretRepository {
	return &secretRepository{
		baseRepository: newBaseRepository[models.Secret](db),
	}
}

// GetByWorkspaceID 根据工作空间ID获取凭证列表
func (r *secretRepository) GetByWorkspaceID(ctx context.Context, workspaceID string, offset, limit int) ([]*models.Secret, int64, error) {
	var secrets []*models.Secret
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Secret{}).Where("workspace_id = ?", workspaceID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Offset(offset).Limit(limit).Find(&secrets).Error
	if err != nil {
		return nil, 0, err
	}

	return secrets, total, nil
}

// GetByName 根据工作空间ID和名称获取凭证
func (r *secretRepository) GetByName(ctx context.Context, workspaceID, name string) (*models.Secret, error) {
	var secret models.Secret
	err := r.db.WithContext(ctx).Where("workspace_id = ? AND name = ?", workspaceID, name).First(&secret).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &secret, nil
}
