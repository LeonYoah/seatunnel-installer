//go:build oracle
// +build oracle

package database

import (
	"fmt"

	"github.com/dzwvip/oracle"
	"gorm.io/gorm"
)

// createOracleDialector 创建Oracle方言器
func createOracleDialector(config *Config) (gorm.Dialector, error) {
	// 使用dzwvip/oracle驱动（已安装在go.mod中）
	return createOracleDialectorWithDzwvip(config)
}

// createOracleDialectorWithDzwvip 使用dzwvip/oracle驱动创建Oracle方言器
func createOracleDialectorWithDzwvip(config *Config) (gorm.Dialector, error) {
	// 构建Oracle DSN
	var dsn string
	if config.ServiceName != "" {
		dsn = fmt.Sprintf("%s/%s@%s:%d/%s",
			config.Username, config.Password, config.Host, config.Port, config.ServiceName)
	} else if config.SID != "" {
		dsn = fmt.Sprintf("%s/%s@%s:%d:%s",
			config.Username, config.Password, config.Host, config.Port, config.SID)
	} else {
		dsn = fmt.Sprintf("%s/%s@%s:%d/%s",
			config.Username, config.Password, config.Host, config.Port, config.Database)
	}

	// 使用dzwvip/oracle驱动创建方言器
	return oracle.Open(dsn), nil
}
