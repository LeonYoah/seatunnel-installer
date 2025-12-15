package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// baseRepository 基础Repository实现
type baseRepository[T any] struct {
	db *gorm.DB
}

// newBaseRepository 创建基础Repository
func newBaseRepository[T any](db *gorm.DB) *baseRepository[T] {
	return &baseRepository[T]{db: db}
}

// Create 创建实体
func (r *baseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// GetByID 根据ID获取实体
func (r *baseRepository[T]) GetByID(ctx context.Context, id string) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

// Update 更新实体
func (r *baseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete 软删除实体
func (r *baseRepository[T]) Delete(ctx context.Context, id string) error {
	var entity T
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity).Error
}

// List 分页查询实体列表
func (r *baseRepository[T]) List(ctx context.Context, offset, limit int) ([]*T, int64, error) {
	var entities []*T
	var total int64

	// 获取总数
	var model T
	if err := r.db.WithContext(ctx).Model(&model).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}
