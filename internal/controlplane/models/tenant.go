package models

import (
	"time"

	"gorm.io/gorm"
)

// Tenant 租户模型
type Tenant struct {
	ID          string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name        string         `gorm:"not null;uniqueIndex;type:varchar(100)" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Status      string         `gorm:"not null;default:'active';type:varchar(20)" json:"status"` // active/inactive
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Workspaces []Workspace `gorm:"foreignKey:TenantID" json:"workspaces,omitempty"`
}

// TableName 指定表名
func (Tenant) TableName() string {
	return "tenants"
}

// BeforeCreate 创建前钩子，生成UUID
func (t *Tenant) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = generateUUID()
	}
	return nil
}
