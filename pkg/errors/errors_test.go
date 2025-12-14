package errors

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(ErrCodeInvalidParam, "invalid parameter")

	if err.Code != ErrCodeInvalidParam {
		t.Errorf("Expected code %s, got %s", ErrCodeInvalidParam, err.Code)
	}

	if err.Message != "invalid parameter" {
		t.Errorf("Expected message 'invalid parameter', got '%s'", err.Message)
	}

	if err.StackTrace == "" {
		t.Error("Expected stack trace to be captured")
	}
}

func TestNewf(t *testing.T) {
	err := Newf(ErrCodeInvalidParam, "invalid parameter: %s", "username")

	if err.Message != "invalid parameter: username" {
		t.Errorf("Expected formatted message, got '%s'", err.Message)
	}
}

func TestWrap(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := Wrap(originalErr, ErrCodeDatabaseError, "database operation failed")

	if wrappedErr.Code != ErrCodeDatabaseError {
		t.Errorf("Expected code %s, got %s", ErrCodeDatabaseError, wrappedErr.Code)
	}

	if wrappedErr.Message != "database operation failed" {
		t.Errorf("Expected message 'database operation failed', got '%s'", wrappedErr.Message)
	}

	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Expected wrapped error to contain original error")
	}
}

func TestWrapNil(t *testing.T) {
	wrappedErr := Wrap(nil, ErrCodeDatabaseError, "should be nil")

	if wrappedErr != nil {
		t.Error("Expected nil when wrapping nil error")
	}
}

func TestWrapAppError(t *testing.T) {
	originalErr := New(ErrCodeInvalidParam, "original error")
	wrappedErr := Wrap(originalErr, ErrCodeDatabaseError, "wrapped error")

	if wrappedErr.Code != ErrCodeDatabaseError {
		t.Errorf("Expected code %s, got %s", ErrCodeDatabaseError, wrappedErr.Code)
	}

	// 应该能够unwrap到原始AppError
	var appErr *AppError
	if !errors.As(wrappedErr, &appErr) {
		t.Error("Expected to unwrap to AppError")
	}
}

func TestWrapf(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := Wrapf(originalErr, ErrCodeDatabaseError, "database operation failed: %s", "connection timeout")

	if wrappedErr.Message != "database operation failed: connection timeout" {
		t.Errorf("Expected formatted message, got '%s'", wrappedErr.Message)
	}
}

func TestWithContext(t *testing.T) {
	err := New(ErrCodeInvalidParam, "invalid parameter")
	err.WithContext("field", "username").WithContext("value", "")

	field, ok := err.GetContext("field")
	if !ok || field != "username" {
		t.Error("Expected to get context value 'username'")
	}

	value, ok := err.GetContext("value")
	if !ok || value != "" {
		t.Error("Expected to get context value ''")
	}

	_, ok = err.GetContext("nonexistent")
	if ok {
		t.Error("Expected false for nonexistent context key")
	}
}

func TestIs(t *testing.T) {
	err := New(ErrCodeInvalidParam, "invalid parameter")

	if !Is(err, ErrCodeInvalidParam) {
		t.Error("Expected Is to return true for matching code")
	}

	if Is(err, ErrCodeDatabaseError) {
		t.Error("Expected Is to return false for non-matching code")
	}

	standardErr := errors.New("standard error")
	if Is(standardErr, ErrCodeInvalidParam) {
		t.Error("Expected Is to return false for standard error")
	}
}

func TestGetCode(t *testing.T) {
	err := New(ErrCodeInvalidParam, "invalid parameter")

	code := GetCode(err)
	if code != ErrCodeInvalidParam {
		t.Errorf("Expected code %s, got %s", ErrCodeInvalidParam, code)
	}

	standardErr := errors.New("standard error")
	code = GetCode(standardErr)
	if code != "" {
		t.Errorf("Expected empty code for standard error, got %s", code)
	}
}

func TestGetMessage(t *testing.T) {
	err := New(ErrCodeInvalidParam, "invalid parameter")

	message := GetMessage(err)
	if message != "invalid parameter" {
		t.Errorf("Expected message 'invalid parameter', got '%s'", message)
	}

	standardErr := errors.New("standard error")
	message = GetMessage(standardErr)
	if message != "standard error" {
		t.Errorf("Expected message 'standard error', got '%s'", message)
	}
}

func TestGetStackTrace(t *testing.T) {
	err := New(ErrCodeInvalidParam, "invalid parameter")

	stack := GetStackTrace(err)
	if stack == "" {
		t.Error("Expected non-empty stack trace")
	}

	if !strings.Contains(stack, "errors_test.go") {
		t.Error("Expected stack trace to contain test file name")
	}

	standardErr := errors.New("standard error")
	stack = GetStackTrace(standardErr)
	if stack != "" {
		t.Error("Expected empty stack trace for standard error")
	}
}

func TestErrorString(t *testing.T) {
	err := New(ErrCodeInvalidParam, "invalid parameter")
	errStr := err.Error()

	if !strings.Contains(errStr, string(ErrCodeInvalidParam)) {
		t.Error("Expected error string to contain error code")
	}

	if !strings.Contains(errStr, "invalid parameter") {
		t.Error("Expected error string to contain message")
	}

	wrappedErr := Wrap(errors.New("original"), ErrCodeDatabaseError, "wrapped")
	errStr = wrappedErr.Error()

	if !strings.Contains(errStr, "original") {
		t.Error("Expected error string to contain original error")
	}
}

func TestIsClientError(t *testing.T) {
	clientErr := New(ErrCodeInvalidParam, "client error")
	if !IsClientError(clientErr) {
		t.Error("Expected IsClientError to return true for client error")
	}

	serverErr := New(ErrCodeDatabaseError, "server error")
	if IsClientError(serverErr) {
		t.Error("Expected IsClientError to return false for server error")
	}
}

func TestIsServerError(t *testing.T) {
	serverErr := New(ErrCodeDatabaseError, "server error")
	if !IsServerError(serverErr) {
		t.Error("Expected IsServerError to return true for server error")
	}

	clientErr := New(ErrCodeInvalidParam, "client error")
	if IsServerError(clientErr) {
		t.Error("Expected IsServerError to return false for client error")
	}
}

func TestIsBusinessError(t *testing.T) {
	businessErr := New(ErrCodeClusterUnavailable, "business error")
	if !IsBusinessError(businessErr) {
		t.Error("Expected IsBusinessError to return true for business error")
	}

	clientErr := New(ErrCodeInvalidParam, "client error")
	if IsBusinessError(clientErr) {
		t.Error("Expected IsBusinessError to return false for client error")
	}
}

func TestUnwrap(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := Wrap(originalErr, ErrCodeDatabaseError, "wrapped error")

	unwrapped := errors.Unwrap(wrappedErr)
	if unwrapped != originalErr {
		t.Error("Expected Unwrap to return original error")
	}
}

func TestErrorChain(t *testing.T) {
	// 创建错误链
	err1 := errors.New("base error")
	err2 := Wrap(err1, ErrCodeDatabaseError, "database error")
	err3 := Wrap(err2, ErrCodeInternalError, "internal error")

	// 测试errors.Is
	if !errors.Is(err3, err1) {
		t.Error("Expected errors.Is to find base error in chain")
	}

	// 测试errors.As
	var appErr *AppError
	if !errors.As(err3, &appErr) {
		t.Error("Expected errors.As to find AppError in chain")
	}

	if appErr.Code != ErrCodeInternalError {
		t.Errorf("Expected code %s, got %s", ErrCodeInternalError, appErr.Code)
	}
}

func TestContextNil(t *testing.T) {
	err := &AppError{
		Code:    ErrCodeInvalidParam,
		Message: "test",
	}

	_, ok := err.GetContext("key")
	if ok {
		t.Error("Expected false for nil context")
	}

	err.WithContext("key", "value")
	val, ok := err.GetContext("key")
	if !ok || val != "value" {
		t.Error("Expected to get context value after WithContext")
	}
}

// 基准测试
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New(ErrCodeInvalidParam, "test error")
	}
}

func BenchmarkWrap(b *testing.B) {
	err := errors.New("original error")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Wrap(err, ErrCodeDatabaseError, "wrapped error")
	}
}

func BenchmarkWithContext(b *testing.B) {
	err := New(ErrCodeInvalidParam, "test error")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err.WithContext("key", "value")
	}
}

// 示例测试
func ExampleNew() {
	err := New(ErrCodeInvalidParam, "username is required")
	fmt.Println(err.Code)
	fmt.Println(err.Message)
	// Output:
	// 1001
	// username is required
}

func ExampleWrap() {
	originalErr := errors.New("connection refused")
	err := Wrap(originalErr, ErrCodeDatabaseError, "failed to connect to database")
	fmt.Println(GetCode(err))
	fmt.Println(GetMessage(err))
	// Output:
	// 2001
	// failed to connect to database
}

func ExampleAppError_WithContext() {
	err := New(ErrCodeInvalidParam, "validation failed")
	err.WithContext("field", "email").WithContext("value", "invalid@")

	field, _ := err.GetContext("field")
	fmt.Println(field)
	// Output:
	// email
}
