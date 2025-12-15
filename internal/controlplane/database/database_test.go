package database

import (
	"os"
	"testing"
)

func TestDatabaseManager(t *testing.T) {
	// 创建临时SQLite数据库进行测试
	config := &Config{
		Type:       "sqlite",
		SQLiteFile: "test_manager.db",
	}

	// 创建数据库管理器
	manager, err := NewManager(config)
	if err != nil {
		t.Fatalf("创建数据库管理器失败: %v", err)
	}
	defer func() {
		manager.Close()
		os.Remove("test_manager.db")
	}()

	// 测试数据库连接
	if err := manager.Ping(); err != nil {
		t.Fatalf("数据库连接测试失败: %v", err)
	}

	// 测试获取DB实例
	db := manager.DB()
	if db == nil {
		t.Fatal("获取数据库实例失败")
	}

	// 测试获取Repository管理器
	repoMgr := manager.Repository()
	if repoMgr == nil {
		t.Fatal("获取Repository管理器失败")
	}

	t.Log("数据库管理器测试通过")
}

func TestDatabaseConfig(t *testing.T) {
	// 测试默认配置
	config := DefaultConfig()
	if config.Type != "sqlite" {
		t.Fatalf("默认数据库类型应该是sqlite, 实际: %s", config.Type)
	}

	// 测试配置验证
	if err := config.Validate(); err != nil {
		t.Fatalf("默认配置验证失败: %v", err)
	}

	// 测试MySQL配置
	mysqlConfig := &Config{
		Type:     "mysql",
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}

	expectedDSN := "root:password@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	if mysqlConfig.DSN() != expectedDSN {
		t.Fatalf("MySQL DSN不匹配: 期望 %s, 实际 %s", expectedDSN, mysqlConfig.DSN())
	}

	// 测试PostgreSQL配置
	pgConfig := &Config{
		Type:     "postgres",
		Host:     "localhost",
		Port:     5432,
		Database: "test",
		Username: "postgres",
		Password: "password",
		SSLMode:  "disable",
	}

	expectedPgDSN := "host=localhost user=postgres password=password dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	if pgConfig.DSN() != expectedPgDSN {
		t.Fatalf("PostgreSQL DSN不匹配: 期望 %s, 实际 %s", expectedPgDSN, pgConfig.DSN())
	}

	t.Log("数据库配置测试通过")
}

func TestOracleConfig(t *testing.T) {
	// 测试Oracle默认配置
	oracleConfig := DefaultOracleConfig()
	if oracleConfig.Type != "oracle" {
		t.Fatalf("Oracle默认数据库类型应该是oracle, 实际: %s", oracleConfig.Type)
	}

	// 测试Oracle配置验证
	if err := oracleConfig.Validate(); err != nil {
		t.Fatalf("Oracle默认配置验证失败: %v", err)
	}

	// 测试Oracle ServiceName配置
	oracleServiceConfig := &Config{
		Type:        "oracle",
		Host:        "localhost",
		Port:        1521,
		Username:    "system",
		Password:    "password",
		ServiceName: "XE",
	}

	expectedOracleDSN := "system/password@localhost:1521/XE"
	if oracleServiceConfig.DSN() != expectedOracleDSN {
		t.Fatalf("Oracle ServiceName DSN不匹配: 期望 %s, 实际 %s", expectedOracleDSN, oracleServiceConfig.DSN())
	}

	// 测试Oracle SID配置
	oracleSIDConfig := &Config{
		Type:     "oracle",
		Host:     "localhost",
		Port:     1521,
		Username: "system",
		Password: "password",
		SID:      "ORCL",
	}

	expectedOracleSIDDSN := "system/password@localhost:1521:ORCL"
	if oracleSIDConfig.DSN() != expectedOracleSIDDSN {
		t.Fatalf("Oracle SID DSN不匹配: 期望 %s, 实际 %s", expectedOracleSIDDSN, oracleSIDConfig.DSN())
	}

	// 测试Oracle Database配置（当ServiceName和SID都为空时）
	oracleDatabaseConfig := &Config{
		Type:     "oracle",
		Host:     "localhost",
		Port:     1521,
		Username: "system",
		Password: "password",
		Database: "XE",
	}

	expectedOracleDbDSN := "system/password@localhost:1521/XE"
	if oracleDatabaseConfig.DSN() != expectedOracleDbDSN {
		t.Fatalf("Oracle Database DSN不匹配: 期望 %s, 实际 %s", expectedOracleDbDSN, oracleDatabaseConfig.DSN())
	}

	// 测试Oracle配置验证 - 缺少必要参数
	invalidOracleConfig := &Config{
		Type: "oracle",
		Host: "localhost",
		Port: 1521,
		// 缺少Username
	}

	if err := invalidOracleConfig.Validate(); err == nil {
		t.Fatal("Oracle配置验证应该失败，因为缺少用户名")
	}

	// 测试Oracle配置验证 - 缺少服务标识
	invalidOracleConfig2 := &Config{
		Type:     "oracle",
		Host:     "localhost",
		Port:     1521,
		Username: "system",
		Password: "password",
		// 缺少ServiceName、SID和Database
	}

	if err := invalidOracleConfig2.Validate(); err == nil {
		t.Fatal("Oracle配置验证应该失败，因为缺少ServiceName、SID或Database")
	}

	t.Log("Oracle配置测试通过")
}

func TestOracleConnection(t *testing.T) {
	// 这个测试只验证Oracle配置和DSN生成，不实际连接数据库
	// 要进行实际连接测试，需要安装Oracle驱动并有可用的Oracle实例

	oracleConfig := &Config{
		Type:        "oracle",
		Host:        "localhost",
		Port:        1521,
		Username:    "system",
		Password:    "password",
		ServiceName: "XE",
	}

	// 测试配置验证
	if err := oracleConfig.Validate(); err != nil {
		t.Fatalf("Oracle配置验证失败: %v", err)
	}

	// 测试DSN生成
	expectedDSN := "system/password@localhost:1521/XE"
	if oracleConfig.DSN() != expectedDSN {
		t.Fatalf("Oracle DSN不匹配: 期望 %s, 实际 %s", expectedDSN, oracleConfig.DSN())
	}

	// 尝试创建连接（会尝试连接Oracle数据库）
	_, err := NewConnection(oracleConfig)
	if err == nil {
		t.Log("Oracle连接成功创建（Oracle数据库可用）")
	} else {
		// 预期会失败，因为没有实际的Oracle数据库实例
		t.Logf("Oracle连接失败（预期行为，无可用Oracle实例）: %v", err)
	}

	t.Log("Oracle连接测试完成")
}

func TestMain(m *testing.M) {
	// 测试前清理
	os.Remove("test.db")

	// 运行测试
	code := m.Run()

	// 测试后清理
	os.Remove("test.db")

	os.Exit(code)
}
func TestAllDatabaseTypes(t *testing.T) {
	// 测试所有支持的数据库类型的配置验证

	testCases := []struct {
		name   string
		config *Config
		valid  bool
	}{
		{
			name: "SQLite配置",
			config: &Config{
				Type:       "sqlite",
				SQLiteFile: "test.db",
			},
			valid: true,
		},
		{
			name: "MySQL配置",
			config: &Config{
				Type:     "mysql",
				Host:     "localhost",
				Port:     3306,
				Database: "test",
				Username: "root",
				Password: "password",
			},
			valid: true,
		},
		{
			name: "PostgreSQL配置",
			config: &Config{
				Type:     "postgres",
				Host:     "localhost",
				Port:     5432,
				Database: "test",
				Username: "postgres",
				Password: "password",
				SSLMode:  "disable",
			},
			valid: true,
		},
		{
			name: "Oracle ServiceName配置",
			config: &Config{
				Type:        "oracle",
				Host:        "localhost",
				Port:        1521,
				Username:    "system",
				Password:    "password",
				ServiceName: "XE",
			},
			valid: true,
		},
		{
			name: "Oracle SID配置",
			config: &Config{
				Type:     "oracle",
				Host:     "localhost",
				Port:     1521,
				Username: "system",
				Password: "password",
				SID:      "ORCL",
			},
			valid: true,
		},
		{
			name: "不支持的数据库类型",
			config: &Config{
				Type: "unsupported",
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			if tc.valid && err != nil {
				t.Errorf("配置应该有效，但验证失败: %v", err)
			}
			if !tc.valid && err == nil {
				t.Errorf("配置应该无效，但验证通过")
			}

			// 测试DSN生成（仅对有效配置）
			if tc.valid && tc.config.Type != "unsupported" {
				dsn := tc.config.DSN()
				if dsn == "" {
					t.Errorf("DSN不应该为空")
				}
				t.Logf("%s DSN: %s", tc.name, dsn)
			}
		})
	}
}
