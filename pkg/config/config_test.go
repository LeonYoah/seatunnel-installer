package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// 创建临时配置文件
	configContent := `
server:
  host: "127.0.0.1"
  port: 9090

database:
  type: "oracle"
  host: "oracle-server"
  port: 1521
  username: "test_user"
  password: "test_password"
  service_name: "TEST"
  max_open_conns: 50
  max_idle_conns: 5

logger:
  level: "debug"
  output_paths:
    - "stdout"
`

	// 写入临时文件
	tmpFile := "test_config.yaml"
	err := os.WriteFile(tmpFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("创建临时配置文件失败: %v", err)
	}
	defer os.Remove(tmpFile)

	// 加载配置
	config, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	// 验证服务器配置
	if config.Server.Host != "127.0.0.1" {
		t.Errorf("期望服务器主机为 '127.0.0.1', 实际: %s", config.Server.Host)
	}
	if config.Server.Port != 9090 {
		t.Errorf("期望服务器端口为 9090, 实际: %d", config.Server.Port)
	}

	// 验证数据库配置
	if config.Database.Type != "oracle" {
		t.Errorf("期望数据库类型为 'oracle', 实际: %s", config.Database.Type)
	}
	if config.Database.Host != "oracle-server" {
		t.Errorf("期望数据库主机为 'oracle-server', 实际: %s", config.Database.Host)
	}
	if config.Database.Port != 1521 {
		t.Errorf("期望数据库端口为 1521, 实际: %d", config.Database.Port)
	}
	if config.Database.ServiceName != "TEST" {
		t.Errorf("期望Oracle服务名为 'TEST', 实际: %s", config.Database.ServiceName)
	}
	if config.Database.MaxOpenConns != 50 {
		t.Errorf("期望最大连接数为 50, 实际: %d", config.Database.MaxOpenConns)
	}

	// 验证日志配置
	if config.Logger.Level != "debug" {
		t.Errorf("期望日志级别为 'debug', 实际: %s", config.Logger.Level)
	}

	t.Log("配置加载测试通过")
}

func TestToDatabaseConfig(t *testing.T) {
	// 创建应用配置
	appConfig := &Config{
		Database: DatabaseConfig{
			Type:            "oracle",
			Host:            "localhost",
			Port:            1521,
			Username:        "system",
			Password:        "password",
			ServiceName:     "XE",
			MaxOpenConns:    50,
			MaxIdleConns:    5,
			ConnMaxLifetime: 7200,
			ConnMaxIdleTime: 900,
		},
	}

	// 转换为数据库配置
	dbConfig := appConfig.ToDatabaseConfig()

	// 验证转换结果
	if dbConfig.Type != "oracle" {
		t.Errorf("期望数据库类型为 'oracle', 实际: %s", dbConfig.Type)
	}
	if dbConfig.Host != "localhost" {
		t.Errorf("期望数据库主机为 'localhost', 实际: %s", dbConfig.Host)
	}
	if dbConfig.Port != 1521 {
		t.Errorf("期望数据库端口为 1521, 实际: %d", dbConfig.Port)
	}
	if dbConfig.ServiceName != "XE" {
		t.Errorf("期望Oracle服务名为 'XE', 实际: %s", dbConfig.ServiceName)
	}
	if dbConfig.MaxOpenConns != 50 {
		t.Errorf("期望最大连接数为 50, 实际: %d", dbConfig.MaxOpenConns)
	}

	t.Log("配置转换测试通过")
}

func TestConfigValidation(t *testing.T) {
	testCases := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "有效的SQLite配置",
			config: &Config{
				Server: ServerConfig{
					Port: 8080,
				},
				Database: DatabaseConfig{
					Type: "sqlite",
				},
				Logger: LoggerConfig{
					Level: "info",
				},
			},
			wantErr: false,
		},
		{
			name: "有效的Oracle配置",
			config: &Config{
				Server: ServerConfig{
					Port: 8080,
				},
				Database: DatabaseConfig{
					Type: "oracle",
				},
				Logger: LoggerConfig{
					Level: "info",
				},
			},
			wantErr: false,
		},
		{
			name: "无效的数据库类型",
			config: &Config{
				Server: ServerConfig{
					Port: 8080,
				},
				Database: DatabaseConfig{
					Type: "unsupported",
				},
				Logger: LoggerConfig{
					Level: "info",
				},
			},
			wantErr: true,
		},
		{
			name: "无效的端口",
			config: &Config{
				Server: ServerConfig{
					Port: 70000,
				},
				Database: DatabaseConfig{
					Type: "sqlite",
				},
				Logger: LoggerConfig{
					Level: "info",
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validate(tc.config)
			if tc.wantErr && err == nil {
				t.Errorf("期望验证失败，但验证通过")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("期望验证通过，但验证失败: %v", err)
			}
		})
	}
}
