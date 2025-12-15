# SeaTunnel企业平台配置指南

## 配置文件位置

SeaTunnel企业平台使用YAML格式的配置文件，支持以下位置：

1. **项目根目录**: `config.yaml`
2. **系统配置目录**: `/etc/seatunnel/config.yaml`
3. **用户配置目录**: `$HOME/.seatunnel/config.yaml`
4. **自定义路径**: 通过命令行参数 `--config` 指定

## 配置文件结构

```yaml
server:
  host: "0.0.0.0"
  port: 8080

database:
  type: "sqlite"  # sqlite, mysql, postgres, oracle
  # ... 数据库特定配置

logger:
  level: "info"
  output_paths:
    - "stdout"
```

## 数据库配置

### SQLite配置（开发环境推荐）

```yaml
database:
  type: "sqlite"
  sqlite_file: "data/seatunnel.db"
  max_open_conns: 1
  max_idle_conns: 1
```

### MySQL配置（生产环境推荐）

```yaml
database:
  type: "mysql"
  host: "localhost"
  port: 3306
  database: "seatunnel"
  username: "seatunnel_user"
  password: "your_password"
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 3600
  conn_max_idle_time: 1800
```

### PostgreSQL配置（生产环境推荐）

```yaml
database:
  type: "postgres"
  host: "localhost"
  port: 5432
  database: "seatunnel"
  username: "seatunnel_user"
  password: "your_password"
  ssl_mode: "disable"
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 3600
  conn_max_idle_time: 1800
```

### Oracle配置（企业级环境）

```yaml
database:
  type: "oracle"
  host: "localhost"
  port: 1521
  username: "system"
  password: "your_password"
  
  # 连接方式1: ServiceName（推荐）
  service_name: "XE"
  
  # 连接方式2: SID（可选）
  # sid: "ORCL"
  
  # 连接方式3: Database（备用）
  # database: "XE"
  
  max_open_conns: 50
  max_idle_conns: 5
  conn_max_lifetime: 7200
  conn_max_idle_time: 900
```

## 配置参数说明

### 服务器配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `server.host` | string | "0.0.0.0" | 服务器监听地址 |
| `server.port` | int | 8080 | 服务器监听端口 |

### 数据库配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `database.type` | string | "sqlite" | 数据库类型 |
| `database.host` | string | "localhost" | 数据库主机地址 |
| `database.port` | int | - | 数据库端口 |
| `database.database` | string | - | 数据库名称 |
| `database.username` | string | - | 数据库用户名 |
| `database.password` | string | - | 数据库密码 |
| `database.sqlite_file` | string | - | SQLite文件路径 |
| `database.service_name` | string | - | Oracle服务名 |
| `database.sid` | string | - | Oracle SID |
| `database.ssl_mode` | string | "disable" | SSL模式 |
| `database.max_open_conns` | int | 100 | 最大打开连接数 |
| `database.max_idle_conns` | int | 10 | 最大空闲连接数 |
| `database.conn_max_lifetime` | int | 3600 | 连接最大生存时间（秒） |
| `database.conn_max_idle_time` | int | 1800 | 连接最大空闲时间（秒） |

### 日志配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `logger.level` | string | "info" | 日志级别 |
| `logger.output_paths` | []string | ["stdout"] | 日志输出路径 |

## 环境变量

所有配置参数都可以通过环境变量覆盖，环境变量名格式为 `SEATUNNEL_<配置路径>`：

```bash
# 设置数据库类型
export SEATUNNEL_DATABASE_TYPE=oracle

# 设置数据库主机
export SEATUNNEL_DATABASE_HOST=oracle-server.example.com

# 设置数据库端口
export SEATUNNEL_DATABASE_PORT=1521

# 设置Oracle服务名
export SEATUNNEL_DATABASE_SERVICE_NAME=PROD
```

## 配置示例

项目提供了多个配置示例文件：

- `config-examples/config-sqlite.yaml` - SQLite配置示例
- `config-examples/config-mysql.yaml` - MySQL配置示例
- `config-examples/config-postgresql.yaml` - PostgreSQL配置示例
- `config-examples/config-oracle.yaml` - Oracle配置示例

## 使用方法

### 1. 复制配置模板

```bash
# 使用SQLite（开发环境）
cp config-examples/config-sqlite.yaml config.yaml

# 使用MySQL（生产环境）
cp config-examples/config-mysql.yaml config.yaml

# 使用Oracle（企业环境）
cp config-examples/config-oracle.yaml config.yaml
```

### 2. 修改配置

编辑 `config.yaml` 文件，修改数据库连接信息：

```yaml
database:
  type: "oracle"
  host: "your-oracle-server.com"
  port: 1521
  username: "your_username"
  password: "your_password"
  service_name: "your_service_name"
```

### 3. 启动应用

```bash
# 使用默认配置文件
./seatunnel-control-plane

# 使用自定义配置文件
./seatunnel-control-plane --config /path/to/your/config.yaml
```

## 配置验证

应用启动时会自动验证配置文件：

- 检查必需参数是否存在
- 验证参数值的有效性
- 测试数据库连接

如果配置有误，应用会输出详细的错误信息并退出。

## 故障排除

### 1. 配置文件未找到

```
Error: failed to read config file: Config File "config" Not Found
```

**解决方案**：确保配置文件存在于正确的位置，或使用 `--config` 参数指定路径。

### 2. 数据库连接失败

```
Error: failed to connect to database: dial tcp: connect: connection refused
```

**解决方案**：
- 检查数据库服务是否运行
- 验证主机地址和端口是否正确
- 确认网络连接和防火墙设置

### 3. Oracle特定错误

```
Error: ORA-00000: DPI-1047: Cannot locate a 64-bit Oracle Client library
```

**解决方案**：安装Oracle客户端库或使用Oracle Instant Client。

## 最佳实践

1. **生产环境**：使用MySQL或PostgreSQL
2. **开发环境**：使用SQLite
3. **企业环境**：使用Oracle
4. **安全性**：不要在配置文件中硬编码密码，使用环境变量
5. **连接池**：根据应用负载调整连接池参数
6. **日志**：生产环境建议将日志输出到文件