package errors

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
)

func TestRecover(t *testing.T) {
	var recovered interface{}
	var stack []byte

	handler := func(r interface{}, s []byte) {
		recovered = r
		stack = s
	}

	func() {
		defer Recover(handler)
		panic("test panic")
	}()

	if recovered == nil {
		t.Error("Expected panic to be recovered")
	}

	if recovered != "test panic" {
		t.Errorf("Expected recovered value 'test panic', got %v", recovered)
	}

	if len(stack) == 0 {
		t.Error("Expected stack trace to be captured")
	}
}

func TestRecoverNoPanic(t *testing.T) {
	handlerCalled := false

	handler := func(r interface{}, s []byte) {
		handlerCalled = true
	}

	func() {
		defer Recover(handler)
		// 正常执行，不panic
	}()

	if handlerCalled {
		t.Error("Expected handler not to be called when no panic occurs")
	}
}

func TestRecoverWithError(t *testing.T) {
	var err error

	func() {
		defer func() {
			err = RecoverWithError(recover())
		}()
		panic("test panic")
	}()

	if err == nil {
		t.Fatal("Expected error to be returned")
	}

	if !strings.Contains(err.Error(), "test panic") {
		t.Errorf("Expected error to contain panic message, got: %v", err)
	}

	// 检查是否为AppError
	var appErr *AppError
	if !errors.As(err, &appErr) {
		t.Error("Expected error to be AppError")
	}

	if appErr.Code != ErrCodeInternalError {
		t.Errorf("Expected code %s, got %s", ErrCodeInternalError, appErr.Code)
	}

	// 检查stack是否在context中
	_, ok := appErr.GetContext("stack")
	if !ok {
		t.Error("Expected stack trace in context")
	}
}

func TestRecoverWithErrorFromError(t *testing.T) {
	originalErr := errors.New("original error")
	var err error

	func() {
		defer func() {
			err = RecoverWithError(recover())
		}()
		panic(originalErr)
	}()

	if err == nil {
		t.Fatal("Expected error to be returned")
	}

	// 应该包装原始错误
	if !errors.Is(err, originalErr) {
		t.Error("Expected error to wrap original error")
	}
}

func TestRecoverWithErrorNoPanic(t *testing.T) {
	err := RecoverWithError(nil)

	if err != nil {
		t.Errorf("Expected nil error when no panic, got: %v", err)
	}
}

func TestSafeGo(t *testing.T) {
	var wg sync.WaitGroup
	var recovered interface{}

	handler := func(r interface{}, s []byte) {
		recovered = r
		wg.Done()
	}

	wg.Add(1)
	SafeGo(func() {
		panic("goroutine panic")
	}, handler)

	wg.Wait()

	if recovered != "goroutine panic" {
		t.Errorf("Expected recovered value 'goroutine panic', got %v", recovered)
	}
}

func TestSafeGoNoPanic(t *testing.T) {
	var wg sync.WaitGroup
	handlerCalled := false
	executed := false

	handler := func(r interface{}, s []byte) {
		handlerCalled = true
	}

	wg.Add(1)
	SafeGo(func() {
		executed = true
		wg.Done()
	}, handler)

	wg.Wait()

	if !executed {
		t.Error("Expected function to be executed")
	}

	if handlerCalled {
		t.Error("Expected handler not to be called when no panic")
	}
}

func TestSafeGoWithCleanup(t *testing.T) {
	var wg sync.WaitGroup
	var recovered interface{}
	cleanupCalled := false

	handler := func(r interface{}, s []byte) {
		recovered = r
	}

	cleanup := func() {
		cleanupCalled = true
		wg.Done()
	}

	wg.Add(1)
	SafeGoWithCleanup(func() {
		panic("goroutine panic")
	}, cleanup, handler)

	wg.Wait()

	if recovered != "goroutine panic" {
		t.Errorf("Expected recovered value 'goroutine panic', got %v", recovered)
	}

	if !cleanupCalled {
		t.Error("Expected cleanup to be called")
	}
}

func TestSafeGoWithCleanupNoPanic(t *testing.T) {
	var wg sync.WaitGroup
	cleanupCalled := false
	executed := false

	cleanup := func() {
		cleanupCalled = true
		wg.Done()
	}

	wg.Add(1)
	SafeGoWithCleanup(func() {
		executed = true
	}, cleanup, nil)

	wg.Wait()

	if !executed {
		t.Error("Expected function to be executed")
	}

	if !cleanupCalled {
		t.Error("Expected cleanup to be called even without panic")
	}
}

func TestSafeGoWithCleanupPanic(t *testing.T) {
	var wg sync.WaitGroup
	cleanupCalled := false

	cleanup := func() {
		cleanupCalled = true
		wg.Done()
	}

	wg.Add(1)
	SafeGoWithCleanup(func() {
		panic("goroutine panic")
	}, cleanup, nil)

	wg.Wait()

	if !cleanupCalled {
		t.Error("Expected cleanup to be called even when function panics")
	}
}

func TestRecoverMiddleware(t *testing.T) {
	var recovered interface{}

	handler := func(r interface{}, s []byte) {
		recovered = r
	}

	middleware := RecoverMiddleware(func() {
		panic("middleware panic")
	}, handler)

	middleware()

	if recovered != "middleware panic" {
		t.Errorf("Expected recovered value 'middleware panic', got %v", recovered)
	}
}

func TestTryWithCleanup(t *testing.T) {
	cleanupCalled := false

	cleanup := func() {
		cleanupCalled = true
	}

	err := TryWithCleanup(func() error {
		return nil
	}, cleanup)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !cleanupCalled {
		t.Error("Expected cleanup to be called")
	}
}

func TestTryWithCleanupError(t *testing.T) {
	cleanupCalled := false
	originalErr := errors.New("function error")

	cleanup := func() {
		cleanupCalled = true
	}

	err := TryWithCleanup(func() error {
		return originalErr
	}, cleanup)

	if err != originalErr {
		t.Errorf("Expected original error, got: %v", err)
	}

	if !cleanupCalled {
		t.Error("Expected cleanup to be called")
	}
}

func TestTryWithCleanupPanic(t *testing.T) {
	cleanupCalled := false

	cleanup := func() {
		cleanupCalled = true
	}

	err := TryWithCleanup(func() error {
		panic("function panic")
	}, cleanup)

	if err == nil {
		t.Fatal("Expected error from panic")
	}

	if !strings.Contains(err.Error(), "function panic") {
		t.Errorf("Expected error to contain panic message, got: %v", err)
	}

	if !cleanupCalled {
		t.Error("Expected cleanup to be called")
	}
}

func TestTryWithCleanupBothPanic(t *testing.T) {
	cleanupCalled := false

	cleanup := func() {
		cleanupCalled = true
		panic("cleanup panic")
	}

	err := TryWithCleanup(func() error {
		panic("function panic")
	}, cleanup)

	if err == nil {
		t.Fatal("Expected error from panic")
	}

	// 应该包含cleanup的错误信息
	if !strings.Contains(err.Error(), "cleanup failed") {
		t.Errorf("Expected error to mention cleanup failure, got: %v", err)
	}

	if !cleanupCalled {
		t.Error("Expected cleanup to be called")
	}
}

func TestMust(t *testing.T) {
	// 测试成功情况
	value := Must(42, nil)
	if value != 42 {
		t.Errorf("Expected value 42, got %d", value)
	}

	// 测试panic情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected Must to panic on error")
		}
	}()

	_ = Must(0, errors.New("test error"))
}

func TestMustNoError(t *testing.T) {
	// 测试成功情况
	MustNoError(nil)

	// 测试panic情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected MustNoError to panic on error")
		}
	}()

	MustNoError(errors.New("test error"))
}

// 基准测试
func BenchmarkRecover(b *testing.B) {
	handler := func(r interface{}, s []byte) {}

	for i := 0; i < b.N; i++ {
		func() {
			defer Recover(handler)
			if i%2 == 0 {
				panic("test")
			}
		}()
	}
}

func BenchmarkRecoverWithError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() (err error) {
			defer func() {
				err = RecoverWithError(recover())
			}()
			if i%2 == 0 {
				panic("test")
			}
			return nil
		}()
	}
}

func BenchmarkTryWithCleanup(b *testing.B) {
	cleanup := func() {}

	for i := 0; i < b.N; i++ {
		_ = TryWithCleanup(func() error {
			if i%2 == 0 {
				panic("test")
			}
			return nil
		}, cleanup)
	}
}

// 示例测试
func ExampleRecover() {
	func() {
		defer Recover(func(recovered interface{}, stack []byte) {
			fmt.Printf("Recovered: %v\n", recovered)
		})
		panic("something went wrong")
	}()
	// Output:
	// Recovered: something went wrong
}

func ExampleSafeGo() {
	var wg sync.WaitGroup
	wg.Add(1)

	SafeGo(func() {
		defer wg.Done()
		// 可能panic的代码
		panic("goroutine error")
	}, func(recovered interface{}, stack []byte) {
		fmt.Printf("Goroutine panic: %v\n", recovered)
	})

	wg.Wait()
	// Output:
	// Goroutine panic: goroutine error
}

func ExampleTryWithCleanup() {
	err := TryWithCleanup(func() error {
		// 可能panic或返回错误的代码
		panic("operation failed")
	}, func() {
		// 清理资源
		fmt.Println("Cleanup executed")
	})

	if err != nil {
		fmt.Println("Error occurred")
	}
	// Output:
	// Cleanup executed
	// Error occurred
}

func ExampleMust() {
	// 成功情况
	value := Must(42, nil)
	fmt.Println(value)
	// Output:
	// 42
}
