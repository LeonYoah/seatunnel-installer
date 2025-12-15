package models

import "gorm.io/gorm"

// AllModels 返回所有需要迁移的模型
func AllModels() []interface{} {
	return []interface{}{
		&Tenant{},
		&Workspace{},
		&User{},
		&Role{},
		&Permission{},
		&UserRole{},
		&RolePermission{},
		&Host{},
		&Cluster{},
		&Node{},
		&Task{},
		&Run{},
		&AuditLog{},
		&Secret{},
	}
}

// AutoMigrate 执行数据库迁移
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(AllModels()...)
}
