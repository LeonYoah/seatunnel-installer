package repository

import (
	"context"
	"errors"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"gorm.io/gorm"
)

// tenantRepository 租户Repository实现
type tenantRepository struct {
	*baseRepository[models.Tenant]
}

// newTenantRepository 创建租户Repository
func newTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{
		baseRepository: newBaseRepository[models.Tenant](db),
	}
}

// GetByName 根据名称获取租户
func (r *tenantRepository) GetByName(ctx context.Context, name string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&tenant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tenant, nil
}

// ListActive 获取活跃租户列表
func (r *tenantRepository) ListActive(ctx context.Context, offset, limit int) ([]*models.Tenant, int64, error) {
	var tenants []*models.Tenant
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Tenant{}).Where("status = ?", "active").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := r.db.WithContext(ctx).Where("status = ?", "active").Offset(offset).Limit(limit).Find(&tenants).Error
	if err != nil {
		return nil, 0, err
	}

	return tenants, total, nil
}
