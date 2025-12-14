package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// ErrorCode 定义错误码类型
type ErrorCode string

const (
	// 1xxx - 客户端错误
	ErrCodeInvalidParam     ErrorCode = "1001" // 参数验证失败
	ErrCodeAuthFailed       ErrorCode = "1002" // 认证失败
	ErrCodePermissionDenied ErrorCode = "1003" // 权限不足
	ErrCodeResourceNotFound ErrorCode = "1004" // 资源不存在

	// 2xxx - 服务器错误
	ErrCodeDatabaseError   ErrorCode = "2001" // 数据库错误
	ErrCodeExternalService ErrorCode = "2002" // 外部服务调用失败
	ErrCodeInternalError   ErrorCode = "2003" // 内部服务错误

	// 3xxx - 业务错误
	ErrCodeClusterUnavailable  ErrorCode = "3001" // 集群不可用
	ErrCodeInvalidTaskConfig   ErrorCode = "3002" // 任务配置无效
	ErrCodeNodeOffline         ErrorCode = "3003" // 节点离线
	ErrCodeInstallFailed       ErrorCode = "3004" // 安装失败
	ErrCodeSSHConnectionFailed ErrorCode = "3005" // SSH连接失败
	ErrCodeCommandFailed       ErrorCode = "3006" // 命令执行失败
)

// AppError 应用错误结构体
type AppError struct {
	Code       ErrorCode              // 错误码
	Message    string                 // 错误消息
	Err        error                  // 原始错误
	Context    map[string]interface{} // 错误上下文
	StackTrace string                 // 堆栈跟踪
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 实现errors.Unwrap接口
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithContext 添加上下文信息
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// GetContext 获取上下文信息
func (e *AppError) GetContext(key string) (interface{}, bool) {
	if e.Context == nil {
		return nil, false
	}
	val, ok := e.Context[key]
	return val, ok
}

// New 创建新的应用错误
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Context:    make(map[string]interface{}),
		StackTrace: captureStackTrace(2),
	}
}

// Newf 创建新的应用错误（格式化消息）
func Newf(code ErrorCode, format string, args ...interface{}) *AppError {
	return &AppError{
		Code:       code,
		Message:    fmt.Sprintf(format, args...),
		Context:    make(map[string]interface{}),
		StackTrace: captureStackTrace(2),
	}
}

// Wrap 包装现有错误
func Wrap(err error, code ErrorCode, message string) *AppError {
	if err == nil {
		return nil
	}

	// 如果已经是AppError，保留原有信息
	var appErr *AppError
	if errors.As(err, &appErr) {
		return &AppError{
			Code:       code,
			Message:    message,
			Err:        appErr,
			Context:    make(map[string]interface{}),
			StackTrace: captureStackTrace(2),
		}
	}

	return &AppError{
		Code:       code,
		Message:    message,
		Err:        err,
		Context:    make(map[string]interface{}),
		StackTrace: captureStackTrace(2),
	}
}

// Wrapf 包装现有错误（格式化消息）
func Wrapf(err error, code ErrorCode, format string, args ...interface{}) *AppError {
	if err == nil {
		return nil
	}

	// 如果已经是AppError，保留原有信息
	var appErr *AppError
	if errors.As(err, &appErr) {
		return &AppError{
			Code:       code,
			Message:    fmt.Sprintf(format, args...),
			Err:        appErr,
			Context:    make(map[string]interface{}),
			StackTrace: captureStackTrace(2),
		}
	}

	return &AppError{
		Code:       code,
		Message:    fmt.Sprintf(format, args...),
		Err:        err,
		Context:    make(map[string]interface{}),
		StackTrace: captureStackTrace(2),
	}
}

// Is 检查错误是否匹配指定的错误码
func Is(err error, code ErrorCode) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

// GetCode 获取错误码
func GetCode(err error) ErrorCode {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return ""
}

// GetMessage 获取错误消息
func GetMessage(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Message
	}
	return err.Error()
}

// GetStackTrace 获取堆栈跟踪
func GetStackTrace(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.StackTrace
	}
	return ""
}

// captureStackTrace 捕获堆栈跟踪
func captureStackTrace(skip int) string {
	const maxDepth = 32
	var pcs [maxDepth]uintptr
	n := runtime.Callers(skip, pcs[:])

	var builder strings.Builder
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		builder.WriteString(fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	return builder.String()
}

// IsClientError 判断是否为客户端错误（1xxx）
func IsClientError(err error) bool {
	code := GetCode(err)
	return len(code) > 0 && code[0] == '1'
}

// IsServerError 判断是否为服务器错误（2xxx）
func IsServerError(err error) bool {
	code := GetCode(err)
	return len(code) > 0 && code[0] == '2'
}

// IsBusinessError 判断是否为业务错误（3xxx）
func IsBusinessError(err error) bool {
	code := GetCode(err)
	return len(code) > 0 && code[0] == '3'
}
