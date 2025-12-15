package models

import (
	"time"

	"gorm.io/gorm"
)

// AuditLog 审计日志模型
type AuditLog struct {
	ID         string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	TenantID   string    `gorm:"not null;type:varchar(36);index" json:"tenant_id"`
	UserID     string    `gorm:"not null;type:varchar(36);index" json:"user_id"`
	Action     string    `gorm:"not null;type:varchar(50)" json:"action"`   // 操作类型：create/update/delete/execute等
	Resource   string    `gorm:"not null;type:varchar(50)" json:"resource"` // 资源类型：cluster/task/host等
	ResourceID string    `gorm:"type:varchar(36)" json:"resource_id"`       // 资源ID
	Details    string    `gorm:"type:longtext" json:"details"`              // 操作详情（JSON）
	Result     string    `gorm:"not null;type:varchar(20)" json:"result"`   // 结果：success/failure
	ErrorMsg   string    `gorm:"type:text" json:"error_msg"`                // 错误信息
	CreatedAt  time.Time `json:"created_at"`
	// 注意：审计日志不允许更新和删除，所以没有UpdatedAt和DeletedAt
}

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_logs"
}

// BeforeCreate 创建前钩子，生成UUID
func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = generateUUID()
	}
	return nil
}
