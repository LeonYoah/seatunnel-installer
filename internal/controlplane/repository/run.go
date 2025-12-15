package repository

import (
	"context"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"gorm.io/gorm"
)

// runRepository 运行Repository实现
type runRepository struct {
	*baseRepository[models.Run]
}

// newRunRepository 创建运行Repository
func newRunRepository(db *gorm.DB) RunRepository {
	return &runRepository{
		baseRepository: newBaseRepository[models.Run](db),
	}
}

// GetByTaskID 根据任务ID获取运行列表
func (r *runRepository) GetByTaskID(ctx context.Context, taskID string, offset, limit int) ([]*models.Run, int64, error) {
	var runs []*models.Run
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Run{}).Where("task_id = ?", taskID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&runs).Error
	if err != nil {
		return nil, 0, err
	}

	return runs, total, nil
}

// GetByClusterID 根据集群ID获取运行列表
func (r *runRepository) GetByClusterID(ctx context.Context, clusterID string, offset, limit int) ([]*models.Run, int64, error) {
	var runs []*models.Run
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Run{}).Where("cluster_id = ?", clusterID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := r.db.WithContext(ctx).Where("cluster_id = ?", clusterID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&runs).Error
	if err != nil {
		return nil, 0, err
	}

	return runs, total, nil
}

// UpdateStatus 更新运行状态
func (r *runRepository) UpdateStatus(ctx context.Context, id, status string) error {
	return r.db.WithContext(ctx).Model(&models.Run{}).Where("id = ?", id).Update("status", status).Error
}

// GetRunningRuns 获取正在运行的任务
func (r *runRepository) GetRunningRuns(ctx context.Context) ([]*models.Run, error) {
	var runs []*models.Run
	err := r.db.WithContext(ctx).Where("status IN ?", []string{"pending", "running"}).Find(&runs).Error
	return runs, err
}
