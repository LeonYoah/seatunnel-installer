package models

import (
	"time"

	"gorm.io/gorm"
)

// Run 运行模型
type Run struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	TaskID    string         `gorm:"not null;type:varchar(36);index" json:"task_id"`
	ClusterID string         `gorm:"not null;type:varchar(36);index" json:"cluster_id"`
	Status    string         `gorm:"not null;default:'pending';type:varchar(20)" json:"status"` // pending/running/succeeded/failed/stopped
	StartTime *time.Time     `json:"start_time"`
	EndTime   *time.Time     `json:"end_time"`
	Duration  int64          `gorm:"default:0" json:"duration"` // 执行时间（秒）
	ErrorMsg  string         `gorm:"type:text" json:"error_msg"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Task    Task    `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	Cluster Cluster `gorm:"foreignKey:ClusterID" json:"cluster,omitempty"`
}

// TableName 指定表名
func (Run) TableName() string {
	return "runs"
}

// BeforeCreate 创建前钩子，生成UUID
func (r *Run) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = generateUUID()
	}
	return nil
}
