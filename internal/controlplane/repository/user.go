package repository

import (
	"context"
	"time"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"gorm.io/gorm"
)

// userRepository 用户Repository实现
type userRepository struct {
	*baseRepository[models.User]
}

// NewUserRepository 创建用户Repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		baseRepository: newBaseRepository[models.User](db),
	}
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByTenantID 根据租户ID获取用户列表
func (r *userRepository) GetByTenantID(ctx context.Context, tenantID string, offset, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	// 获取总数
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("tenant_id = ?", tenantID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).
		Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetByWorkspaceID 根据工作空间ID获取用户列表
func (r *userRepository) GetByWorkspaceID(ctx context.Context, workspaceID string, offset, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	// 获取总数
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("workspace_id = ?", workspaceID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).
		Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateLastLogin 更新最后登录时间
func (r *userRepository) UpdateLastLogin(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).
		Update("last_login_at", &now).Error
}

// GetWithRoles 获取用户及其角色信息
func (r *userRepository) GetWithRoles(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Preload("Roles").Preload("Roles.Permissions").
		Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// AssignRole 为用户分配角色
func (r *userRepository) AssignRole(ctx context.Context, userID, roleID string) error {
	userRole := models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.WithContext(ctx).Create(&userRole).Error
}

// RemoveRole 移除用户角色
func (r *userRepository) RemoveRole(ctx context.Context, userID, roleID string) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&models.UserRole{}).Error
}
