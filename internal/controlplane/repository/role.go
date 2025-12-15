package repository

import (
	"context"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"gorm.io/gorm"
)

// roleRepository 角色Repository实现
type roleRepository struct {
	*baseRepository[models.Role]
}

// NewRoleRepository 创建角色Repository
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		baseRepository: newBaseRepository[models.Role](db),
	}
}

// GetByTenantID 根据租户ID获取角色列表
func (r *roleRepository) GetByTenantID(ctx context.Context, tenantID string, offset, limit int) ([]*models.Role, int64, error) {
	var roles []*models.Role
	var total int64

	// 获取总数
	err := r.db.WithContext(ctx).Model(&models.Role{}).Where("tenant_id = ?", tenantID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).
		Offset(offset).Limit(limit).Find(&roles).Error
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// GetByName 根据租户ID和角色名获取角色
func (r *roleRepository) GetByName(ctx context.Context, tenantID, name string) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND name = ?", tenantID, name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetWithPermissions 获取角色及其权限信息
func (r *roleRepository) GetWithPermissions(ctx context.Context, id string) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).Preload("Permissions").Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// AssignPermission 为角色分配权限
func (r *roleRepository) AssignPermission(ctx context.Context, roleID, permissionID string) error {
	rolePermission := models.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	return r.db.WithContext(ctx).Create(&rolePermission).Error
}

// RemovePermission 移除角色权限
func (r *roleRepository) RemovePermission(ctx context.Context, roleID, permissionID string) error {
	return r.db.WithContext(ctx).Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&models.RolePermission{}).Error
}

// GetBuiltInRoles 获取内置角色
func (r *roleRepository) GetBuiltInRoles(ctx context.Context) ([]*models.Role, error) {
	var roles []*models.Role
	err := r.db.WithContext(ctx).Where("is_built_in = ?", true).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// permissionRepository 权限Repository实现
type permissionRepository struct {
	*baseRepository[models.Permission]
}

// NewPermissionRepository 创建权限Repository
func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{
		baseRepository: newBaseRepository[models.Permission](db),
	}
}

// GetByPermission 根据权限字符串获取权限
func (r *permissionRepository) GetByPermission(ctx context.Context, permission string) (*models.Permission, error) {
	var perm models.Permission
	err := r.db.WithContext(ctx).Where("permission = ?", permission).First(&perm).Error
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

// GetByResource 根据资源获取权限列表
func (r *permissionRepository) GetByResource(ctx context.Context, resource string) ([]*models.Permission, error) {
	var permissions []*models.Permission
	err := r.db.WithContext(ctx).Where("resource = ?", resource).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// ListAll 获取所有权限
func (r *permissionRepository) ListAll(ctx context.Context) ([]*models.Permission, error) {
	var permissions []*models.Permission
	err := r.db.WithContext(ctx).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}
