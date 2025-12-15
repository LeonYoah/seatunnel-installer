//go:build !oracle
// +build !oracle

package database

import (
	"errors"

	"gorm.io/gorm"
)

// createOracleDialector 创建Oracle方言器（存根版本）
func createOracleDialector(config *Config) (gorm.Dialector, error) {
	return nil, errors.New("Oracle支持未启用。要启用Oracle支持，请使用 -tags oracle 构建标签")
}
