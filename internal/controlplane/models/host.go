package models

import (
	"time"

	"gorm.io/gorm"
)

// Host 主机模型
type Host struct {
	ID            string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	WorkspaceID   string         `gorm:"not null;type:varchar(36);index" json:"workspace_id"`
	Name          string         `gorm:"not null;type:varchar(100)" json:"name"`
	IP            string         `gorm:"not null;type:varchar(45)" json:"ip"`
	Port          int            `gorm:"not null;default:22" json:"port"`
	User          string         `gorm:"type:varchar(50)" json:"user"`
	AuthType      string         `gorm:"type:varchar(20);default:'password'" json:"auth_type"` // password/key
	Password      string         `gorm:"type:text" json:"-"`                                   // 加密存储，不返回给前端
	KeyPath       string         `gorm:"type:varchar(500)" json:"key_path"`
	Description   string         `gorm:"type:text" json:"description"`
	AgentStatus   string         `gorm:"not null;default:'not-installed';type:varchar(20)" json:"agent_status"` // not-installed/installed/offline
	Status        string         `gorm:"not null;default:'offline';type:varchar(20)" json:"status"`             // online/offline
	CPU           int            `gorm:"default:0" json:"cpu"`
	Memory        int64          `gorm:"default:0" json:"memory"`
	Disk          int64          `gorm:"default:0" json:"disk"`
	CPUUsage      float64        `gorm:"default:0" json:"cpu_usage"`
	MemoryUsage   float64        `gorm:"default:0" json:"memory_usage"`
	LastHeartbeat *time.Time     `json:"last_heartbeat"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Workspace Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
}

// TableName 指定表名
func (Host) TableName() string {
	return "hosts"
}

// BeforeCreate 创建前钩子，生成UUID
func (h *Host) BeforeCreate(tx *gorm.DB) error {
	if h.ID == "" {
		h.ID = generateUUID()
	}
	return nil
}
