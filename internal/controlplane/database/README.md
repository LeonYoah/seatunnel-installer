# 数据库支持

SeaTunnel企业平台支持多种数据库作为后端存储，包括SQLite、MySQL、PostgreSQL和Oracle。

## 支持的数据库

| 数据库 | 状态 | 驱动 | 说明 |
|--------|------|------|------|
| SQLite | ✅ 完全支持 | gorm.io/driver/sqlite | 默认数据库，适合开发和小规模部署 |
| MySQL | ✅ 完全支持 | gorm.io/driver/mysql | 生产环境推荐 |
| PostgreSQL | ✅ 完全支持 | gorm.io/driver/postgres | 生产环境推荐 |
| Oracle | ✅ 完全支持 | github.com/dzwvip/oracle | 企业级数据库支持，已默认启用 |

## 文件结构

```
internal/controlplane/database/
├── config.go          # 数据库配置管理
├── connection.go      # 数据库连接管理
├── manager.go         # 数据库管理器
├── sqlite.go          # SQLite支持（需要CGO）
├── sqlite_nocgo.go    # SQLite无CGO支持
├── oracle.go          # Oracle数据库支持
├── database_test.go   # 数据库测试
└── README.md          # 本文档
```

## Oracle数据库适配

### 配置支持

Oracle数据库配置支持以下参数：

```go
type Config struct {
    Type        string `json:"type"`         // "oracle"
    Host        string `json:"host"`         // Oracle主机地址
    Port        int    `json:"port"`         // Oracle端口（默认1521）
    Username    string `json:"username"`     // 数据库用户名
    Password    string `json:"password"`     // 数据库密码
    Database    string `json:"database"`     // 数据库名
    ServiceName string `json:"service_name"` // Oracle服务名
    SID         string `json:"sid"`          // Oracle SID
    // ... 其他配置
}
```

### DSN生成

系统支持三种Oracle连接方式：

1. **ServiceName方式**（推荐）：`user/password@host:port/service_name`
2. **SID方式**：`user/password@host:port:sid`
3. **Database方式**：`user/password@host:port/database`

### 驱动依赖

Oracle支持使用go-ora驱动，这是一个纯Go实现，无需Oracle客户端：

```bash
go get github.com/sijms/go-ora/v2
go get github.com/dzwvip/oracle
```
