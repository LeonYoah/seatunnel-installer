package main

import (
	"fmt"
	"os"

	"github.com/seatunnel/enterprise-platform/pkg/errors"
)

func main() {
	fmt.Println("=== 错误处理示例 ===\n")

	// 示例1: 创建和包装错误
	example1()

	// 示例2: 错误上下文
	example2()

	// 示例3: Panic恢复
	example3()

	// 示例4: 安全的Goroutine
	example4()

	// 示例5: Try-Cleanup模式
	example5()
}

// 示例1: 创建和包装错误
func example1() {
	fmt.Println("示例1: 创建和包装错误")

	// 创建新错误
	err := errors.New(errors.ErrCodeInvalidParam, "用户名不能为空")
	fmt.Printf("错误: %v\n", err)
	fmt.Printf("错误码: %s\n", errors.GetCode(err))
	fmt.Printf("错误消息: %s\n", errors.GetMessage(err))

	// 包装错误
	wrappedErr := errors.Wrap(err, errors.ErrCodeDatabaseError, "保存用户失败")
	fmt.Printf("包装后的错误: %v\n", wrappedErr)

	// 检查错误类型
	if errors.IsClientError(err) {
		fmt.Println("这是一个客户端错误")
	}

	fmt.Println()
}

// 示例2: 错误上下文
func example2() {
	fmt.Println("示例2: 错误上下文")

	err := errors.New(errors.ErrCodeInvalidParam, "验证失败")
	err.WithContext("field", "email").
		WithContext("value", "invalid@").
		WithContext("rule", "email格式")

	fmt.Printf("错误: %v\n", err)

	if field, ok := err.GetContext("field"); ok {
		fmt.Printf("字段: %v\n", field)
	}
	if value, ok := err.GetContext("value"); ok {
		fmt.Printf("值: %v\n", value)
	}
	if rule, ok := err.GetContext("rule"); ok {
		fmt.Printf("规则: %v\n", rule)
	}

	fmt.Println()
}

// 示例3: Panic恢复
func example3() {
	fmt.Println("示例3: Panic恢复")

	// 使用Recover
	func() {
		defer errors.Recover(func(recovered interface{}, stack []byte) {
			fmt.Printf("捕获到panic: %v\n", recovered)
			// 在实际应用中，这里应该记录日志
		})

		// 模拟panic
		panic("出错了！")
	}()

	// 使用RecoverWithError
	err := func() (err error) {
		defer func() {
			err = errors.RecoverWithError(recover())
		}()

		// 模拟panic
		panic("另一个错误")
	}()

	if err != nil {
		fmt.Printf("从panic恢复的错误: %v\n", err)
	}

	fmt.Println()
}

// 示例4: 安全的Goroutine
func example4() {
	fmt.Println("示例4: 安全的Goroutine")

	// 使用SafeGo
	done := make(chan bool)
	errors.SafeGo(func() {
		defer func() { done <- true }()
		// 模拟可能panic的操作
		panic("goroutine中的错误")
	}, func(recovered interface{}, stack []byte) {
		fmt.Printf("Goroutine panic: %v\n", recovered)
	})
	<-done

	// 使用SafeGoWithCleanup
	done2 := make(chan bool)
	errors.SafeGoWithCleanup(
		func() {
			// 模拟操作
			panic("需要清理的错误")
		},
		func() {
			// 清理资源
			fmt.Println("执行清理操作")
			done2 <- true
		},
		func(recovered interface{}, stack []byte) {
			fmt.Printf("Goroutine panic: %v\n", recovered)
		},
	)
	<-done2

	fmt.Println()
}

// 示例5: Try-Cleanup模式
func example5() {
	fmt.Println("示例5: Try-Cleanup模式")

	// 模拟文件操作
	err := errors.TryWithCleanup(
		func() error {
			// 模拟打开文件
			fmt.Println("打开文件...")

			// 模拟操作失败
			return errors.New(errors.ErrCodeInternalError, "文件读取失败")
		},
		func() {
			// 清理资源
			fmt.Println("关闭文件...")
		},
	)

	if err != nil {
		fmt.Printf("操作失败: %v\n", err)
	}

	// 使用Must辅助函数
	fmt.Println("\n使用Must辅助函数:")
	config := errors.Must(loadConfig())
	fmt.Printf("配置: %v\n", config)

	fmt.Println()
}

// 辅助函数
func loadConfig() (string, error) {
	// 模拟加载配置
	return "config-data", nil
}

// 模拟数据库操作
func saveUser(username string) error {
	if username == "" {
		return errors.New(errors.ErrCodeInvalidParam, "用户名不能为空")
	}

	// 模拟数据库错误
	dbErr := fmt.Errorf("connection refused")
	return errors.Wrap(dbErr, errors.ErrCodeDatabaseError, "保存用户失败")
}

// 模拟SSH操作
func connectSSH(host string) error {
	if host == "" {
		return errors.New(errors.ErrCodeInvalidParam, "主机地址不能为空")
	}

	// 模拟连接失败
	return errors.Newf(errors.ErrCodeSSHConnectionFailed, "无法连接到主机: %s", host)
}

// 模拟命令执行
func executeCommand(cmd string) error {
	if cmd == "" {
		return errors.New(errors.ErrCodeInvalidParam, "命令不能为空")
	}

	// 模拟命令执行失败
	return errors.New(errors.ErrCodeCommandFailed, "命令执行失败").
		WithContext("command", cmd).
		WithContext("exit_code", 1)
}

// 演示完整的错误处理流程
func demonstrateFullFlow() {
	fmt.Println("=== 完整错误处理流程 ===\n")

	// 1. 尝试保存用户
	if err := saveUser(""); err != nil {
		handleError(err)
	}

	// 2. 尝试SSH连接
	if err := connectSSH(""); err != nil {
		handleError(err)
	}

	// 3. 尝试执行命令
	if err := executeCommand(""); err != nil {
		handleError(err)
	}
}

// 统一的错误处理函数
func handleError(err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		fmt.Printf("错误码: %s\n", appErr.Code)
		fmt.Printf("错误消息: %s\n", appErr.Message)

		if len(appErr.Context) > 0 {
			fmt.Println("上下文信息:")
			for k, v := range appErr.Context {
				fmt.Printf("  %s: %v\n", k, v)
			}
		}

		// 根据错误类型采取不同的处理策略
		if errors.IsClientError(err) {
			fmt.Println("处理策略: 返回400错误给客户端")
		} else if errors.IsServerError(err) {
			fmt.Println("处理策略: 返回500错误，记录详细日志")
		} else if errors.IsBusinessError(err) {
			fmt.Println("处理策略: 返回业务错误码，提示用户")
		}
	} else {
		fmt.Printf("标准错误: %v\n", err)
	}
	fmt.Println()
}

// 演示在HTTP处理器中使用
func httpHandlerExample() {
	handler := errors.RecoverMiddleware(func() {
		// HTTP处理逻辑
		panic("HTTP处理器panic")
	}, func(recovered interface{}, stack []byte) {
		fmt.Printf("HTTP处理器panic: %v\n", recovered)
		// 返回500错误给客户端
	})

	handler()
}

// 演示在后台任务中使用
func backgroundTaskExample() {
	errors.SafeGoWithCleanup(
		func() {
			// 后台任务逻辑
			fmt.Println("执行后台任务...")

			// 模拟错误
			if err := processTask(); err != nil {
				panic(err)
			}
		},
		func() {
			// 清理资源
			fmt.Println("清理后台任务资源...")
		},
		func(recovered interface{}, stack []byte) {
			fmt.Printf("后台任务失败: %v\n", recovered)
			// 记录日志，发送告警
		},
	)
}

func processTask() error {
	return errors.New(errors.ErrCodeInternalError, "任务处理失败")
}

// 演示文件操作的错误处理
func fileOperationExample() error {
	var file *os.File
	return errors.TryWithCleanup(
		func() error {
			// 打开文件
			var err error
			file, err = os.Open("config.yaml")
			if err != nil {
				return errors.Wrap(err, errors.ErrCodeInternalError, "打开配置文件失败")
			}

			// 读取文件
			// ... 处理文件内容

			return nil
		},
		func() {
			// 确保文件被关闭
			if file != nil {
				file.Close()
				fmt.Println("关闭文件")
			}
		},
	)
}
