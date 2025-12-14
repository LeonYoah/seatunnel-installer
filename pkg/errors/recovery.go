package errors

import (
	"fmt"
	"runtime/debug"
)

// RecoveryHandler 定义panic恢复处理函数类型
type RecoveryHandler func(recovered interface{}, stack []byte)

// DefaultRecoveryHandler 默认的panic恢复处理函数
var DefaultRecoveryHandler RecoveryHandler = func(recovered interface{}, stack []byte) {
	fmt.Printf("Panic recovered: %v\nStack trace:\n%s\n", recovered, string(stack))
}

// Recover 恢复panic并调用处理函数
// 使用示例:
//
//	defer errors.Recover(func(recovered interface{}, stack []byte) {
//	    logger.Error("Panic recovered", zap.Any("panic", recovered), zap.String("stack", string(stack)))
//	})
func Recover(handler RecoveryHandler) {
	if r := recover(); r != nil {
		stack := debug.Stack()
		if handler != nil {
			handler(r, stack)
		} else {
			DefaultRecoveryHandler(r, stack)
		}
	}
}

// RecoverWithError 恢复panic并返回错误
// 使用示例:
//
//	func someFunc() (err error) {
//	    defer func() {
//	        err = errors.RecoverWithError(recover())
//	    }()
//	    // ... 可能panic的代码
//	    return nil
//	}
func RecoverWithError(recovered interface{}) error {
	if recovered == nil {
		return nil
	}

	stack := debug.Stack()

	// 如果recovered是error类型，包装它
	if err, ok := recovered.(error); ok {
		return Wrap(err, ErrCodeInternalError, fmt.Sprintf("panic recovered: %v", recovered)).
			WithContext("stack", string(stack))
	}

	// 否则创建新错误
	return New(ErrCodeInternalError, fmt.Sprintf("panic recovered: %v", recovered)).
		WithContext("stack", string(stack))
}

// SafeGo 安全地启动goroutine，自动恢复panic
// 使用示例:
//
//	errors.SafeGo(func() {
//	    // ... 可能panic的代码
//	}, func(recovered interface{}, stack []byte) {
//	    logger.Error("Goroutine panic", zap.Any("panic", recovered))
//	})
func SafeGo(fn func(), handler RecoveryHandler) {
	go func() {
		defer Recover(handler)
		fn()
	}()
}

// SafeGoWithCleanup 安全地启动goroutine，支持清理函数
// 使用示例:
//
//	errors.SafeGoWithCleanup(
//	    func() {
//	        // ... 可能panic的代码
//	    },
//	    func() {
//	        // ... 清理资源
//	    },
//	    func(recovered interface{}, stack []byte) {
//	        logger.Error("Goroutine panic", zap.Any("panic", recovered))
//	    },
//	)
func SafeGoWithCleanup(fn func(), cleanup func(), handler RecoveryHandler) {
	go func() {
		defer func() {
			// 先恢复panic
			if r := recover(); r != nil {
				stack := debug.Stack()

				// 执行清理
				if cleanup != nil {
					// 清理函数也可能panic，需要再次保护
					defer func() {
						if r2 := recover(); r2 != nil {
							fmt.Printf("Cleanup function panicked: %v\n", r2)
						}
					}()
					cleanup()
				}

				// 调用处理函数
				if handler != nil {
					handler(r, stack)
				} else {
					DefaultRecoveryHandler(r, stack)
				}
			} else if cleanup != nil {
				// 正常退出也执行清理
				cleanup()
			}
		}()
		fn()
	}()
}

// RecoverMiddleware 创建一个恢复中间件（用于HTTP处理器等）
// 使用示例:
//
//	handler := errors.RecoverMiddleware(
//	    func() {
//	        // ... HTTP处理逻辑
//	    },
//	    func(recovered interface{}, stack []byte) {
//	        logger.Error("HTTP handler panic", zap.Any("panic", recovered))
//	    },
//	)
func RecoverMiddleware(handler func(), recoveryHandler RecoveryHandler) func() {
	return func() {
		defer Recover(recoveryHandler)
		handler()
	}
}

// TryWithCleanup 尝试执行函数，如果panic则恢复并执行清理
// 使用示例:
//
//	err := errors.TryWithCleanup(
//	    func() error {
//	        // ... 可能panic或返回错误的代码
//	        return nil
//	    },
//	    func() {
//	        // ... 清理资源
//	    },
//	)
func TryWithCleanup(fn func() error, cleanup func()) (err error) {
	defer func() {
		// 恢复panic
		if r := recover(); r != nil {
			err = RecoverWithError(r)
		}

		// 执行清理（无论是否panic）
		if cleanup != nil {
			defer func() {
				if r := recover(); r != nil {
					// 清理函数panic，包装到错误中
					cleanupErr := RecoverWithError(r)
					if err != nil {
						err = Wrap(err, ErrCodeInternalError, fmt.Sprintf("cleanup failed: %v", cleanupErr))
					} else {
						err = cleanupErr
					}
				}
			}()
			cleanup()
		}
	}()

	return fn()
}

// Must 如果错误不为nil则panic
// 使用示例:
//
//	config := errors.Must(loadConfig())
func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

// MustNoError 如果错误不为nil则panic（无返回值版本）
// 使用示例:
//
//	errors.MustNoError(saveConfig(config))
func MustNoError(err error) {
	if err != nil {
		panic(err)
	}
}
