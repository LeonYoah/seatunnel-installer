//go:build cgo

package database

import (
	"fmt"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// createSQLiteDialector 创建SQLite方言器（需要CGO支持）
func createSQLiteDialector(config *Config) (gorm.Dialector, error) {
	// 确保SQLite文件目录存在
	if err := ensureDir(filepath.Dir(config.SQLiteFile)); err != nil {
		return nil, fmt.Errorf("创建SQLite目录失败: %w", err)
	}
	return sqlite.Open(config.DSN()), nil
}
