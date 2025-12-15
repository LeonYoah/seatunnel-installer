package models

import (
	"time"

	"gorm.io/gorm"
)

// Secret 凭证模型
type Secret struct {
	ID          string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	WorkspaceID string         `gorm:"not null;type:varchar(36);index" json:"workspace_id"`
	Name        string         `gorm:"not null;type:varchar(100)" json:"name"`
	Type        string         `gorm:"not null;type:varchar(20)" json:"type"` // database/api/ssh
	Value       string         `gorm:"not null;type:longtext" json:"-"`       // 加密的凭证值，不返回给前端
	CreatedBy   string         `gorm:"not null;type:varchar(36)" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Workspace Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
}

// TableName 指定表名
func (Secret) TableName() string {
	return "secrets"
}

// BeforeCreate 创建前钩子，生成UUID
func (s *Secret) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = generateUUID()
	}
	return nil
}
