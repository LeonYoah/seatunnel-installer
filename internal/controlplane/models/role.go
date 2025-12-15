package models

import (
	"time"

	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	ID          string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	TenantID    string         `gorm:"not null;type:varchar(36);index" json:"tenant_id"`
	Name        string         `gorm:"not null;type:varchar(50)" json:"name"`
	DisplayName string         `gorm:"not null;type:varchar(100)" json:"display_name"`
	Description string         `gorm:"type:text" json:"description"`
	IsBuiltIn   bool           `gorm:"not null;default:false" json:"is_built_in"`                // 是否为内置角色
	Status      string         `gorm:"not null;default:'active';type:varchar(20)" json:"status"` // active/inactive
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Tenant      Tenant       `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Users       []User       `gorm:"many2many:user_roles;" json:"users,omitempty"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// BeforeCreate 创建前钩子，生成UUID
func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = generateUUID()
	}
	return nil
}

// Permission 权限模型
type Permission struct {
	ID          string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Permission  string         `gorm:"not null;uniqueIndex;type:varchar(100)" json:"permission"` // 格式: resource:action
	Resource    string         `gorm:"not null;type:varchar(50)" json:"resource"`
	Action      string         `gorm:"not null;type:varchar(50)" json:"action"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Roles []Role `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// BeforeCreate 创建前钩子，生成UUID
func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = generateUUID()
	}
	return nil
}

// UserRole 用户角色关联表
type UserRole struct {
	UserID string `gorm:"primaryKey;type:varchar(36)" json:"user_id"`
	RoleID string `gorm:"primaryKey;type:varchar(36)" json:"role_id"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	RoleID       string `gorm:"primaryKey;type:varchar(36)" json:"role_id"`
	PermissionID string `gorm:"primaryKey;type:varchar(36)" json:"permission_id"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}
