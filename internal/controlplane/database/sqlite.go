package database

import (
	"fmt"
	"path/filepath"

	"github.com/glebarez/sqlite" // 纯Go SQLite驱动，专为GORM设计，无需CGO
	"gorm.io/gorm"
)

// createSQLiteDialector 创建SQLite方言器（使用glebarez/sqlite，零依赖）
func createSQLiteDialector(config *Config) (gorm.Dialector, error) {
	// 确保SQLite文件目录存在
	if err := ensureDir(filepath.Dir(config.SQLiteFile)); err != nil {
		return nil, fmt.Errorf("创建SQLite目录失败: %w", err)
	}

	// 使用glebarez/sqlite驱动，专为GORM设计的纯Go SQLite实现
	dsn := config.DSN()

	// 添加SQLite特定的参数以优化性能
	if dsn == config.SQLiteFile {
		dsn = fmt.Sprintf("%s?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=cache_size(-64000)&_pragma=foreign_keys(ON)", dsn)
	}

	return sqlite.Open(dsn), nil
}
