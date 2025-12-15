package models

import (
	"time"

	"gorm.io/gorm"
)

// Workspace 工作空间模型
type Workspace struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	TenantID  string         `gorm:"not null;type:varchar(36);index" json:"tenant_id"`
	Name      string         `gorm:"not null;type:varchar(100)" json:"name"`
	Status    string         `gorm:"not null;default:'active';type:varchar(20)" json:"status"` // active/inactive
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Tenant   Tenant    `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Hosts    []Host    `gorm:"foreignKey:WorkspaceID" json:"hosts,omitempty"`
	Clusters []Cluster `gorm:"foreignKey:WorkspaceID" json:"clusters,omitempty"`
	Tasks    []Task    `gorm:"foreignKey:WorkspaceID" json:"tasks,omitempty"`
	Secrets  []Secret  `gorm:"foreignKey:WorkspaceID" json:"secrets,omitempty"`
}

// TableName 指定表名
func (Workspace) TableName() string {
	return "workspaces"
}

// BeforeCreate 创建前钩子，生成UUID
func (w *Workspace) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = generateUUID()
	}
	return nil
}
