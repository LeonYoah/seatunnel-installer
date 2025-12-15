package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID          string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	TenantID    string         `gorm:"not null;type:varchar(36);index" json:"tenant_id"`
	WorkspaceID string         `gorm:"not null;type:varchar(36);index" json:"workspace_id"`
	Username    string         `gorm:"not null;uniqueIndex;type:varchar(50)" json:"username"`
	Email       string         `gorm:"not null;uniqueIndex;type:varchar(100)" json:"email"`
	Password    string         `gorm:"not null;type:varchar(255)" json:"-"`                      // 不在JSON中返回密码
	Status      string         `gorm:"not null;default:'active';type:varchar(20)" json:"status"` // active/inactive/locked
	LastLoginAt *time.Time     `json:"last_login_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Tenant    Tenant    `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Workspace Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	Roles     []Role    `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前钩子，生成UUID和加密密码
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = generateUUID()
	}
	return u.HashPassword()
}

// BeforeUpdate 更新前钩子，如果密码被修改则重新加密
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 检查密码是否被修改
	if tx.Statement.Changed("Password") {
		return u.HashPassword()
	}
	return nil
}

// HashPassword 加密密码
func (u *User) HashPassword() error {
	if u.Password == "" {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// HasPermission 检查用户是否有指定权限
func (u *User) HasPermission(resource, action string) bool {
	permission := resource + ":" + action
	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			if perm.Permission == permission {
				return true
			}
		}
	}
	return false
}

// HasRole 检查用户是否有指定角色
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.Status == "active"
}
