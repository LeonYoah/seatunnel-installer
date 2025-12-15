package models

import (
	"time"

	"gorm.io/gorm"
)

// Cluster 集群模型
type Cluster struct {
	ID          string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	WorkspaceID string         `gorm:"not null;type:varchar(36);index" json:"workspace_id"`
	Name        string         `gorm:"not null;type:varchar(100)" json:"name"`
	Version     string         `gorm:"not null;type:varchar(50)" json:"version"`
	DeployMode  string         `gorm:"not null;type:varchar(20)" json:"deploy_mode"`                  // hybrid/separated
	DeployType  string         `gorm:"not null;type:varchar(20)" json:"deploy_type"`                  // baremetal/docker/k8s
	Status      string         `gorm:"not null;default:'registering';type:varchar(20)" json:"status"` // registering/ready/unhealthy/offline
	MasterNodes string         `gorm:"type:text" json:"master_nodes"`                                 // JSON数组存储节点ID列表
	WorkerNodes string         `gorm:"type:text" json:"worker_nodes"`                                 // JSON数组存储节点ID列表
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Workspace Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	Nodes     []Node    `gorm:"foreignKey:ClusterID" json:"nodes,omitempty"`
	Tasks     []Task    `gorm:"foreignKey:ClusterID" json:"tasks,omitempty"`
	Runs      []Run     `gorm:"foreignKey:ClusterID" json:"runs,omitempty"`
}

// TableName 指定表名
func (Cluster) TableName() string {
	return "clusters"
}

// BeforeCreate 创建前钩子，生成UUID
func (c *Cluster) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = generateUUID()
	}
	return nil
}
