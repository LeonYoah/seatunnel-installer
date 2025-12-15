package models

import (
	"time"

	"gorm.io/gorm"
)

// Node 节点模型
type Node struct {
	ID            string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	ClusterID     string         `gorm:"not null;type:varchar(36);index" json:"cluster_id"`
	IP            string         `gorm:"not null;type:varchar(45)" json:"ip"`
	Hostname      string         `gorm:"not null;type:varchar(100)" json:"hostname"`
	Role          string         `gorm:"not null;type:varchar(20)" json:"role"` // master/worker
	Version       string         `gorm:"not null;type:varchar(50)" json:"version"`
	Status        string         `gorm:"not null;default:'offline';type:varchar(20)" json:"status"` // online/offline/unhealthy
	CPU           int            `gorm:"default:0" json:"cpu"`
	Memory        int64          `gorm:"default:0" json:"memory"`
	Disk          int64          `gorm:"default:0" json:"disk"`
	LastHeartbeat *time.Time     `json:"last_heartbeat"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Cluster Cluster `gorm:"foreignKey:ClusterID" json:"cluster,omitempty"`
}

// TableName 指定表名
func (Node) TableName() string {
	return "nodes"
}

// BeforeCreate 创建前钩子，生成UUID
func (n *Node) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = generateUUID()
	}
	return nil
}
