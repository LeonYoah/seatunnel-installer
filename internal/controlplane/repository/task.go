package repository

import (
	"context"
	"errors"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"gorm.io/gorm"
)

// taskRepository 任务Repository实现
type taskRepository struct {
	*baseRepository[models.Task]
}

// newTaskRepository 创建任务Repository
func newTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{
		baseRepository: newBaseRepository[models.Task](db),
	}
}

// GetByWorkspaceID 根据工作空间ID获取任务列表
func (r *taskRepository) GetByWorkspaceID(ctx context.Context, workspaceID string, offset, limit int) ([]*models.Task, int64, error) {
	var tasks []*models.Task
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Task{}).Where("workspace_id = ?", workspaceID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Offset(offset).Limit(limit).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// GetByName 根据工作空间ID和名称获取任务
func (r *taskRepository) GetByName(ctx context.Context, workspaceID, name string) (*models.Task, error) {
	var task models.Task
	err := r.db.WithContext(ctx).Where("workspace_id = ? AND name = ?", workspaceID, name).First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

// GetVersions 获取任务的所有版本
func (r *taskRepository) GetVersions(ctx context.Context, taskID string) ([]*models.Task, error) {
	var tasks []*models.Task
	err := r.db.WithContext(ctx).Where("id = ?", taskID).Order("version DESC").Find(&tasks).Error
	return tasks, err
}

// GetLatestVersion 获取任务的最新版本
func (r *taskRepository) GetLatestVersion(ctx context.Context, taskID string) (*models.Task, error) {
	var task models.Task
	err := r.db.WithContext(ctx).Where("id = ?", taskID).Order("version DESC").First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}
