# 错误处理包 (Error Handling Package)

这个包提供了统一的错误处理和恢复机制，用于SeaTunnel企业级运维控制平台。

## 功能特性

### 1. 统一的错误类型定义

- **错误码系统**：使用标准化的错误码（1xxx客户端错误、2xxx服务器错误、3xxx业务错误）
- **错误包装**：支持错误链，保留原始错误信息
- **上下文信息**：可以为错误添加任意上下文数据
- **堆栈跟踪**：自动捕获错误发生时的堆栈信息

### 2. 错误包装和上下文传递

- 支持标准库`errors.Is`和`errors.As`
- 可以在错误传递过程中添加上下文信息
- 保留完整的错误链

### 3. Panic恢复和清理逻辑

- 自动恢复panic
- 支持清理函数
- 安全的goroutine启动
- 中间件模式支持

## 使用示例

### 创建错误

```go
import "github.com/your-org/seatunnel-platform/pkg/errors"

// 创建新错误
err := errors.New(errors.ErrCodeInvalidParam, "用户名不能为空")

// 创建格式化错误
err := errors.Newf(errors.ErrCodeInvalidParam, "用户名长度必须在%d到%d之间", 3, 20)
```

### 包装错误

```go
// 包装标准错误
if err := db.Query(); err != nil {
    return errors.Wrap(err, errors.ErrCodeDatabaseError, "查询用户失败")
}

// 包装并格式化
if err := db.Query(); err != nil {
    return errors.Wrapf(err, errors.ErrCodeDatabaseError, "查询用户失败: %s", username)
}
```

### 添加上下文

```go
err := errors.New(errors.ErrCodeInvalidParam, "验证失败")
err.WithContext("field", "email").WithContext("value", "invalid@")

// 获取上下文
if field, ok := err.GetContext("field"); ok {
    fmt.Printf("字段: %v\n", field)
}
```

### 错误检查

```go
// 检查错误码
if errors.Is(err, errors.ErrCodeInvalidParam) {
    // 处理参数错误
}

// 获取错误码
code := errors.GetCode(err)

// 检查错误类型
if errors.IsClientError(err) {
    // 客户端错误
}
```

### Panic恢复

```go
// 基本恢复
func someFunc() {
    defer errors.Recover(func(recovered interface{}, stack []byte) {
        logger.Error("Panic recovered", 
            zap.Any("panic", recovered),
            zap.String("stack", string(stack)))
    })
    
    // 可能panic的代码
}

// 恢复并返回错误
func someFunc() (err error) {
    defer func() {
        err = errors.RecoverWithError(recover())
    }()
    
    // 可能panic的代码
    return nil
}
```

### 安全的Goroutine

```go
// 基本用法
errors.SafeGo(func() {
    // 可能panic的代码
}, func(recovered interface{}, stack []byte) {
    logger.Error("Goroutine panic", zap.Any("panic", recovered))
})

// 带清理函数
errors.SafeGoWithCleanup(
    func() {
        // 可能panic的代码
    },
    func() {
        // 清理资源
        conn.Close()
    },
    func(recovered interface{}, stack []byte) {
        logger.Error("Goroutine panic", zap.Any("panic", recovered))
    },
)
```

### Try-Cleanup模式

```go
err := errors.TryWithCleanup(
    func() error {
        // 打开资源
        file, err := os.Open("config.yaml")
        if err != nil {
            return err
        }
        
        // 可能panic或返回错误的操作
        return processFile(file)
    },
    func() {
        // 清理资源（无论是否出错都会执行）
        if file != nil {
            file.Close()
        }
    },
)
```

### Must辅助函数

```go
// 有返回值的Must
config := errors.Must(loadConfig())

// 无返回值的Must
errors.MustNoError(saveConfig(config))
```

## 错误码定义

### 客户端错误 (1xxx)

- `1001` - 参数验证失败
- `1002` - 认证失败
- `1003` - 权限不足
- `1004` - 资源不存在

### 服务器错误 (2xxx)

- `2001` - 数据库错误
- `2002` - 外部服务调用失败
- `2003` - 内部服务错误

### 业务错误 (3xxx)

- `3001` - 集群不可用
- `3002` - 任务配置无效
- `3003` - 节点离线
- `3004` - 安装失败
- `3005` - SSH连接失败
- `3006` - 命令执行失败

## 最佳实践

### 1. 错误创建

- 在边界处创建错误（如API入口、数据库层）
- 使用合适的错误码
- 提供清晰的错误消息

### 2. 错误包装

- 在错误传递过程中添加上下文
- 保留原始错误信息
- 不要重复包装相同的信息

### 3. 错误处理

- 在适当的层级处理错误
- 记录足够的上下文信息
- 对用户友好的错误消息

### 4. Panic恢复

- 在goroutine入口处恢复panic
- 在HTTP处理器中恢复panic
- 在关键操作中使用TryWithCleanup

### 5. 日志记录

```go
import (
    "github.com/your-org/seatunnel-platform/pkg/errors"
    "github.com/your-org/seatunnel-platform/pkg/logger"
    "go.uber.org/zap"
)

func handleError(err error) {
    if appErr, ok := err.(*errors.AppError); ok {
        logger.Error("操作失败",
            zap.String("code", string(appErr.Code)),
            zap.String("message", appErr.Message),
            zap.Any("context", appErr.Context),
            zap.String("stack", appErr.StackTrace),
            zap.Error(appErr.Err),
        )
    } else {
        logger.Error("操作失败", zap.Error(err))
    }
}
```

## 与标准库的兼容性

这个包完全兼容Go标准库的`errors`包：

```go
import (
    "errors"
    apperrors "github.com/your-org/seatunnel-platform/pkg/errors"
)

// 使用标准库的Is
if errors.Is(err, apperrors.ErrCodeInvalidParam) {
    // ...
}

// 使用标准库的As
var appErr *apperrors.AppError
if errors.As(err, &appErr) {
    // ...
}

// 使用标准库的Unwrap
unwrapped := errors.Unwrap(err)
```

## 性能考虑

- 错误创建时会捕获堆栈跟踪，有一定性能开销
- 在高频路径上可以考虑使用标准错误
- 错误包装的开销很小
- 上下文添加使用map，查询效率为O(1)

## 测试

运行测试：

```bash
go test ./pkg/errors/...
```

运行基准测试：

```bash
go test -bench=. ./pkg/errors/...
```

## 示例代码

查看`errors_test.go`和`recovery_test.go`中的示例测试了解更多用法。
