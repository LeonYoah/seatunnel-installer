//go:build !cgo

package database

import (
	"fmt"

	"gorm.io/gorm"
)

// createSQLiteDialector 在没有CGO支持时返回错误
func createSQLiteDialector(config *Config) (gorm.Dialector, error) {
	return nil, fmt.Errorf("SQLite支持需要CGO，请使用MySQL或PostgreSQL")
}
