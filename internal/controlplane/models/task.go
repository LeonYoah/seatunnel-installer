package models

import (
	"time"

	"gorm.io/gorm"
)

// Task 任务模型
type Task struct {
	ID          string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	WorkspaceID string         `gorm:"not null;type:varchar(36);index" json:"workspace_id"`
	ClusterID   string         `gorm:"type:varchar(36);index" json:"cluster_id"` // 可选，任务可以不绑定到特定集群
	Name        string         `gorm:"not null;type:varchar(100)" json:"name"`
	Type        string         `gorm:"not null;type:varchar(20)" json:"type"` // batch/stream/dag
	Config      string         `gorm:"type:longtext" json:"config"`           // 任务配置（YAML）
	Version     int            `gorm:"not null;default:1" json:"version"`
	Status      string         `gorm:"not null;default:'draft';type:varchar(20)" json:"status"` // draft/published/archived
	CreatedBy   string         `gorm:"not null;type:varchar(36)" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Workspace Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	Cluster   *Cluster  `gorm:"foreignKey:ClusterID" json:"cluster,omitempty"`
	Runs      []Run     `gorm:"foreignKey:TaskID" json:"runs,omitempty"`
}

// TableName 指定表名
func (Task) TableName() string {
	return "tasks"
}

// BeforeCreate 创建前钩子，生成UUID
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = generateUUID()
	}
	return nil
}
