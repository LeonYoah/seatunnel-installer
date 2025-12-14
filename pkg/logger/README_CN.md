# 日志包

SeaTunnel 企业级平台的生产级日志框架，基于 Uber 的 Zap 日志库构建，支持自动日志轮转。

## 特性

- **多种日志级别**：DEBUG、INFO、WARN、ERROR、FATAL
- **灵活输出**：控制台、文件或同时输出
- **自动日志轮转**：基于文件大小，可配置保留策略
- **结构化日志**：文件输出 JSON 格式，控制台输出人类可读格式
- **线程安全**：支持多协程并发使用
- **零依赖**：仅使用 Zap 和 Lumberjack（业界标准库）

## 安装

日志包是 SeaTunnel 企业级平台的一部分，需要以下依赖：

```bash
go get go.uber.org/zap
go get github.com/natefinch/lumberjack
```

## 快速开始

### 基础用法（仅控制台）

```go
package main

import (
    "github.com/seatunnel/enterprise-platform/pkg/logger"
    "go.uber.org/zap"
)

func main() {
    // 使用默认配置初始化
    config := logger.DefaultConfig()
    logger.Init(config)
    defer logger.Close()

    // 记录日志消息
    logger.Info("应用程序已启动", zap.String("version", "1.0.0"))
    logger.Warn("这是一个警告")
    logger.Error("发生了一个错误", zap.Error(err))
}
```

### 文件输出与轮转

```go
config := &logger.Config{
    Level:      "info",
    OutputPath: "/var/log/seatunnel/app.log",
    MaxSize:    100,  // 每个文件 100 MB
    MaxBackups: 3,    // 保留 3 个旧文件
    MaxAge:     28,   // 保留 28 天
    Compress:   true, // 压缩轮转的文件
    Console:    false,
}

logger.Init(config)
defer logger.Close()

logger.Info("记录到文件并自动轮转")
```

### 同时输出到控制台和文件

```go
config := &logger.Config{
    Level:      "debug",
    OutputPath: "/var/log/seatunnel/app.log",
    MaxSize:    100,
    MaxBackups: 3,
    MaxAge:     28,
    Compress:   true,
    Console:    true, // 同时输出到控制台
}

logger.Init(config)
defer logger.Close()

logger.Debug("这条消息同时出现在控制台和文件中")
```

## 配置

### 配置结构

```go
type Config struct {
    Level      string // 日志级别：debug、info、warn、error、fatal
    OutputPath string // 文件输出路径（空表示仅控制台）
    MaxSize    int    // 轮转前的最大文件大小（MB）
    MaxBackups int    // 保留的旧日志文件数量
    MaxAge     int    // 保留旧日志文件的最大天数
    Compress   bool   // 是否压缩轮转的日志文件
    Console    bool   // 是否输出到控制台
}
```

### 默认配置

```go
&Config{
    Level:      "info",
    OutputPath: "",      // 仅控制台
    MaxSize:    100,     // 100 MB
    MaxBackups: 3,       // 保留 3 个旧文件
    MaxAge:     28,      // 保留 28 天
    Compress:   true,    // 压缩旧文件
    Console:    true,    // 输出到控制台
}
```

## 日志级别

- **DEBUG**：详细的调试信息
- **INFO**：一般性信息消息
- **WARN**：潜在有害情况的警告消息
- **ERROR**：错误事件的错误消息
- **FATAL**：导致应用程序退出的严重错误消息

## 日志轮转

日志轮转由 Lumberjack 自动处理：

- **基于大小**：当文件达到 `MaxSize` MB 时轮转
- **基于时间**：删除超过 `MaxAge` 天的文件
- **基于数量**：仅保留 `MaxBackups` 个旧文件
- **压缩**：可选择使用 gzip 压缩轮转的文件

### 轮转示例

配置 `MaxSize: 100`、`MaxBackups: 3`、`MaxAge: 28` 时：

```
app.log          (当前日志文件)
app.log.1        (最近的备份)
app.log.2        (较旧的备份)
app.log.3        (最旧的备份)
```

当 `app.log` 达到 100MB 时，它会被重命名为 `app.log.1`，并创建新的 `app.log`。

## API 参考

### 初始化函数

- `Init(config *Config) error` - 使用完整配置初始化日志器
- `InitSimple(level string, outputPaths []string) error` - 简单初始化（向后兼容）
- `DefaultConfig() *Config` - 返回默认配置

### 日志记录函数

- `Debug(msg string, fields ...zap.Field)` - 记录调试消息
- `Info(msg string, fields ...zap.Field)` - 记录信息消息
- `Warn(msg string, fields ...zap.Field)` - 记录警告消息
- `Error(msg string, fields ...zap.Field)` - 记录错误消息
- `Fatal(msg string, fields ...zap.Field)` - 记录致命消息并退出

### 工具函数

- `Get() *zap.Logger` - 获取全局日志器实例
- `Sync() error` - 刷新缓冲的日志条目
- `Close() error` - 关闭日志器并释放资源

## 结构化日志

使用 Zap 的字段类型进行结构化日志记录：

```go
logger.Info("用户已登录",
    zap.String("username", "admin"),
    zap.Int("user_id", 123),
    zap.Duration("login_time", time.Since(start)),
)

logger.Error("数据库连接失败",
    zap.Error(err),
    zap.String("host", "localhost"),
    zap.Int("port", 5432),
)
```

## 输出格式

### 控制台输出（人类可读）

```
2025-12-14T20:50:17.188+0800    INFO    logger/logger.go:165    应用程序已启动    {"version": "1.0.0"}
2025-12-14T20:50:17.188+0800    WARN    logger/logger.go:175    警告消息
2025-12-14T20:50:17.188+0800    ERROR   logger/logger.go:170    发生错误         {"error": "连接被拒绝"}
```

### 文件输出（JSON）

```json
{"level":"INFO","time":"2025-12-14T20:50:17.188+0800","caller":"logger/logger.go:165","msg":"应用程序已启动","version":"1.0.0"}
{"level":"WARN","time":"2025-12-14T20:50:17.188+0800","caller":"logger/logger.go:175","msg":"警告消息"}
{"level":"ERROR","time":"2025-12-14T20:50:17.188+0800","caller":"logger/logger.go:170","msg":"发生错误","error":"连接被拒绝"}
```

## 最佳实践

1. **始终关闭**：使用 `defer logger.Close()` 确保资源被释放
2. **初始化一次**：在应用程序启动时初始化日志器一次
3. **使用结构化字段**：优先使用结构化字段而非字符串格式化
4. **选择适当的级别**：开发环境使用 DEBUG，生产环境使用 INFO
5. **配置轮转**：根据磁盘空间设置适当的轮转参数
6. **处理错误**：检查 `Init()`、`Sync()` 和 `Close()` 返回的错误

## 线程安全

日志器完全线程安全，可以从多个协程并发使用，无需额外的同步机制。

## 性能

- **零分配**：Zap 的结构化日志最小化内存分配
- **缓冲 I/O**：日志被缓冲以获得更好的性能
- **异步轮转**：日志轮转不会阻塞日志操作

## 测试

运行测试套件：

```bash
go test -v ./pkg/logger/...
```

## 要求

- Go 1.21 或更高版本
- SeaTunnel 企业级平台规范中的需求 1.2
