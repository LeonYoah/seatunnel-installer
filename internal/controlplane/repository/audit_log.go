package repository

import (
	"context"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"gorm.io/gorm"
)

// auditLogRepository 审计日志Repository实现
type auditLogRepository struct {
	db *gorm.DB
}

// newAuditLogRepository 创建审计日志Repository
func newAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

// Create 创建审计日志
func (r *auditLogRepository) Create(ctx context.Context, log *models.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetByTenantID 根据租户ID获取审计日志列表
func (r *auditLogRepository) GetByTenantID(ctx context.Context, tenantID string, offset, limit int) ([]*models.AuditLog, int64, error) {
	var logs []*models.AuditLog
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.AuditLog{}).Where("tenant_id = ?", tenantID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetByUserID 根据用户ID获取审计日志列表
func (r *auditLogRepository) GetByUserID(ctx context.Context, userID string, offset, limit int) ([]*models.AuditLog, int64, error) {
	var logs []*models.AuditLog
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.AuditLog{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetByResource 根据资源获取审计日志列表
func (r *auditLogRepository) GetByResource(ctx context.Context, tenantID, resource, resourceID string, offset, limit int) ([]*models.AuditLog, int64, error) {
	var logs []*models.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&models.AuditLog{}).Where("tenant_id = ? AND resource = ?", tenantID, resource)
	if resourceID != "" {
		query = query.Where("resource_id = ?", resourceID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
